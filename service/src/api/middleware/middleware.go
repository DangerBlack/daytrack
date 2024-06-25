package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"512b.it/daytrack/src/database"
	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const AUTHORIZATION = "Authorization"
const BEARER = "Bearer"
const KEY = "Key"

type authHeader struct {
	Authorization string `header:"Authorization"`
}

func AuthUserGuards(configuration models.Configuration, db *database.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestCtx := ctx.Request.Context()

		h := authHeader{}

		if err := ctx.ShouldBindHeader(&h); err != nil {
			utils.Logger(requestCtx).Err(err).Msgf("Unable to bind the header")
			ctx.JSON(http.StatusUnauthorized, models.NewError(models.ErrorUnauthorized, ""))
			ctx.Abort()
			return
		}

		if h.Authorization == "" {
			utils.Logger(requestCtx).Error().Msgf("Authorization header is missing")
			ctx.JSON(http.StatusUnauthorized, models.NewError(models.ErrorUnauthorized, ""))
			ctx.Abort()
			return
		}

		if strings.HasPrefix(h.Authorization, BEARER) {
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
				ctx.Set(utils.USER_ID_CONTEXT_KEY, identityID)
			} else {
				utils.Logger(requestCtx).Error().Msgf("Jwt token found is not linked to a user")
				ctx.JSON(http.StatusUnauthorized, models.NewError(models.ErrorUnauthorized, ""))
				ctx.Abort()
				return
			}
		} else if strings.HasPrefix(h.Authorization, KEY) {
			var err error
			var apiKey *models.ApiKey
			key := strings.Split(h.Authorization, KEY+" ")[1]

			if apiKey, err = db.GetAPIKey(key); err != nil {
				utils.Logger(requestCtx).Err(err).Msgf("Unable to get the api key")
				ctx.JSON(http.StatusUnauthorized, models.NewError(models.ErrorUnauthorized, ""))
				ctx.Abort()
				return
			}

			if apiKey == nil {
				utils.Logger(requestCtx).Error().Msgf("Api key not found")
				ctx.JSON(http.StatusUnauthorized, models.NewError(models.ErrorUnauthorized, ""))
				ctx.Abort()
				return
			}

			if apiKey.DeleteAt != nil {
				utils.Logger(requestCtx).Error().Msgf("Api key is deleted")
				ctx.JSON(http.StatusUnauthorized, models.NewError(models.ErrorUnauthorized, ""))
				ctx.Abort()
				return
			}

			ctx.Set(utils.USER_ID_CONTEXT_KEY, strconv.FormatInt(apiKey.UserID, 10))
		} else {
			utils.Logger(requestCtx).Error().Msgf("Authorization header is not a bearer token")
			ctx.JSON(http.StatusUnauthorized, models.NewError(models.ErrorUnauthorized, ""))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func AuthApiKeyGuards(configuration models.Configuration, db *database.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestCtx := ctx.Request.Context()

		var err error
		var apiKey *models.ApiKey
		key := ctx.Query("key")

		if apiKey, err = db.GetAPIKey(key); err != nil {
			utils.Logger(requestCtx).Err(err).Msgf("Unable to get the api key")
			ctx.JSON(http.StatusUnauthorized, models.NewError(models.ErrorUnauthorized, ""))
			ctx.Abort()
			return
		}

		if apiKey == nil {
			utils.Logger(requestCtx).Error().Msgf("Api key not found")
			ctx.JSON(http.StatusUnauthorized, models.NewError(models.ErrorUnauthorized, ""))
			ctx.Abort()
			return
		}

		if apiKey.DeleteAt != nil {
			utils.Logger(requestCtx).Error().Msgf("Api key is deleted")
			ctx.JSON(http.StatusUnauthorized, models.NewError(models.ErrorUnauthorized, ""))
			ctx.Abort()
			return
		}

		ctx.Set(utils.USER_ID_CONTEXT_KEY, fmt.Sprintf("%d", apiKey.UserID))
		ctx.Next()
	}
}
