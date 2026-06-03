package api_key

import (
	"strconv"
	"strings"
	"time"

	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
	"github.com/gin-gonic/gin"
)

type ApiKeyController struct {
	unauthenticatedRoute *gin.RouterGroup
	authenticatedRoute   *gin.RouterGroup

	configuration models.Configuration
	apiKey        *Service
	logger        utils.ContextLogger
}

func Inject(
	unauthenticatedRoute *gin.RouterGroup,
	authenticatedRoute *gin.RouterGroup,
	apiKey *Service,
	configuration models.Configuration,
) {
	controller := new(unauthenticatedRoute, authenticatedRoute, apiKey, configuration)
	controller.injectUnauthenticatedRoutes()
	controller.injectAuthenticatedRoutes()
}

func new(
	unauthenticatedRoute *gin.RouterGroup,
	authenticatedRoute *gin.RouterGroup,
	apiKey *Service,
	configuration models.Configuration,
) *ApiKeyController {
	logger := utils.InitServiceLogger("ApiKeyController")

	return &ApiKeyController{
		unauthenticatedRoute: unauthenticatedRoute,
		authenticatedRoute:   authenticatedRoute,
		apiKey:               apiKey,
		configuration:        configuration,
		logger:               logger,
	}
}

func (c *ApiKeyController) injectUnauthenticatedRoutes() {
}

func (c *ApiKeyController) injectAuthenticatedRoutes() {
	v1 := c.authenticatedRoute.Group("v1")
	{
		v1.POST("/api_keys", c.createApiKeyRoute())
		v1.GET("/api_keys", c.listApiKeysRoute())
		v1.DELETE("/api_keys/:id", c.deleteApiKeyRoute())
	}
}

// @Tags api_key
// @Security TokenAuth
// @Schemes https
// @Router /v1/api_keys [POST]
// @Summary Create an api key
// @Description Create an api key for the authenticated user
// @Accept json
// @Param name query string false "Name of the api key" default(default)
// @Produce json
// @Success 201 {object} models.ApiKey
func (c *ApiKeyController) createApiKeyRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var key *string
		var userID int64

		name := strings.TrimSpace(ctx.DefaultQuery("name", "default"))

		if len(name) < 3 || len(name) > 255 {
			ctx.JSON(400, models.NewError(models.ErrorBadRequest, "name must be between 3 and 255 characters"))
			return
		}

		if !models.IsValidTrackName(name) {
			ctx.JSON(400, models.NewError(models.ErrorBadRequest, "name can only contain letters, numbers, hyphens and underscores"))
			return
		}

		if userID, err = utils.GetAuthenticatedUserID(ctx); err != nil {
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to get authenticated user id"))
			return
		}

		if key, err = c.apiKey.CreateApiKey(userID, name); err != nil {
			c.logger(ctx).Err(err).Msg("Error while creating api key")
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to create api key"))
			return
		}

		ctx.JSON(201, models.ApiKey{
			Name:      name,
			UserID:    userID,
			Key:       *key,
			CreatedAt: time.Now(),
		})
	}
}

// @Tags api_key
// @Security TokenAuth
// @Schemes https
// @Router /v1/api_keys [GET]
// @Summary List api keys
// @Description List all api keys for the authenticated user
// @Accept json
// @Produce json
// @Success 200 {object} models.List[models.ApiKey]
func (c *ApiKeyController) listApiKeysRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var apiKeys []models.ApiKey
		var userID int64

		if userID, err = utils.GetAuthenticatedUserID(ctx); err != nil {
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to get authenticated user id"))
			return
		}

		if apiKeys, err = c.apiKey.ListApiKeys(userID); err != nil {
			c.logger(ctx).Err(err).Msg("Error while listing api keys")
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to list api keys"))
			return
		}

		ctx.JSON(200, models.List[models.ApiKey]{
			Items: apiKeys,
		})
	}
}

// @Tags api_key
// @Security TokenAuth
// @Schemes https
// @Router /v1/api_keys/{id} [DELETE]
// @Summary Delete an api key
// @Description Delete an api key for the authenticated user
// @Accept json
// @Param id path int true "Id of the api key"
// @Produce json
// @Success 204
func (c *ApiKeyController) deleteApiKeyRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var userID int64
		var keyID int64

		if userID, err = utils.GetAuthenticatedUserID(ctx); err != nil {
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to get authenticated user id"))
			return
		}

		keyID, err = strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(400, models.NewError(models.ErrorBadRequest, "invalid api key id"))
			return
		}

		if err = c.apiKey.DeleteApiKey(userID, keyID); err != nil {
			c.logger(ctx).Err(err).Msg("Error while deleting api key")
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to delete api key"))
			return
		}

		ctx.Status(204)
	}
}
