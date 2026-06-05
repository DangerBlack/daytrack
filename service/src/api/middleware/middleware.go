package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"512b.it/daytrack/src/database"
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
			ctx.Set(utils.USER_ID_CONTEXT_KEY, identityID)
		} else {
			utils.Logger(requestCtx).Error().Msgf("Jwt token found is not linked to a user")
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
		key := ctx.Query("key")

		if key == "" {
			username := ctx.Param("username")
			trackName := ctx.Param("track_name")

			if username == "" || trackName == "" {
				ctx.JSON(http.StatusUnauthorized, models.NewError(models.ErrorUnauthorized, ""))
				ctx.Abort()
				return
			}

			user, err := db.GetUserByName(username)
			if err != nil {
				if errors.Is(err, database.ErrorNotFound) {
					ctx.JSON(http.StatusForbidden, models.NewError(models.ErrorForbidden, "event cannot be called by you"))
				} else {
					utils.Logger(requestCtx).Err(err).Msgf("Failed to look up user %s", username)
					ctx.JSON(http.StatusInternalServerError, models.NewError(models.ErrorInternalServerError, ""))
				}
				ctx.Abort()
				return
			}

			track, err := db.GetTrackByUserIDAndName(user.ID, trackName)
			if err != nil {
				if errors.Is(err, database.ErrorNotFound) {
					ctx.JSON(http.StatusForbidden, models.NewError(models.ErrorForbidden, "event cannot be called by you"))
				} else {
					utils.Logger(requestCtx).Err(err).Msgf("Failed to look up track %s for user %d", trackName, user.ID)
					ctx.JSON(http.StatusInternalServerError, models.NewError(models.ErrorInternalServerError, ""))
				}
				ctx.Abort()
				return
			}

			isWrite := ctx.Request.Method == http.MethodPost
			if isWrite && track.Visibility != models.TrackVisibilityPublicWrite {
				ctx.JSON(http.StatusForbidden, models.NewError(models.ErrorForbidden, "event cannot be called by you"))
				ctx.Abort()
				return
			}
			if !isWrite && track.Visibility != models.TrackVisibilityPublicRead && track.Visibility != models.TrackVisibilityPublicWrite {
				ctx.JSON(http.StatusForbidden, models.NewError(models.ErrorForbidden, "event cannot be called by you"))
				ctx.Abort()
				return
			}

			ctx.Next()
			return
		}

		var err error
		var apiKey *models.ApiKey

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
