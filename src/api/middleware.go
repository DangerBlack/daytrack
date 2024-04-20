package api

import (
	"net/http"
	"strings"

	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const AUTHORIZATION = "Authorization"
const BEARER = "Bearer"

type authHeader struct {
	Authorization string `header:"Authorization"`
}

func AuthUserGuards(configuration models.Configuration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestCtx := ctx.Request.Context()

		h := authHeader{}

		if err := ctx.ShouldBindHeader(&h); err != nil {
			utils.Logger(requestCtx).Err(err).Msgf("Unable to bind the header")
			ctx.JSON(http.StatusUnauthorized, models.NewError(models.ErrorUnauthorized, ""))
			ctx.Abort()
			return
		}

		authorizationTokenHeader := strings.Split(h.Authorization, BEARER+" ")

		if len(authorizationTokenHeader) < 2 {
			utils.Logger(requestCtx).Error().Msgf("Token is too short %s", authorizationTokenHeader)
			ctx.JSON(http.StatusUnauthorized, models.NewError(models.ErrorUnauthorized, ""))
			ctx.Abort()
			return
		}

		token, err := models.ParseAndValidateJWT(requestCtx, authorizationTokenHeader[1], configuration.JWT.PublicKey)

		if err != nil {
			utils.Logger(requestCtx).Warn().Msgf("Unable to parse and validate the jwt %s", err)
			ctx.JSON(http.StatusUnauthorized, models.NewError(models.ErrorUnauthorized, ""))
			ctx.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		if identityID, found := claims[models.JWTSubjectClaimKey]; found {
			ctx.Set(models.USER_ID_CONTEXT_KEY, identityID)
		} else {
			utils.Logger(requestCtx).Error().Msgf("Jwt token found is not linked to a user")
			ctx.JSON(http.StatusUnauthorized, models.NewError(models.ErrorUnauthorized, ""))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
