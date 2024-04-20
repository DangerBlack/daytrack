package api_key

import (
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
	}
}

func (c *ApiKeyController) createApiKeyRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var key *string
		var userID int64

		name := ctx.DefaultQuery("name", "default")

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
