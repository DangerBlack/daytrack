package api_key

import (
	"strconv"
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
		//v1.GET("/api_keys", c.getApiKeys)
	}
}

func (c *ApiKeyController) createApiKeyRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var key *string

		userID, _ := strconv.Atoi(ctx.Value(models.USER_ID_CONTEXT_KEY).(string))
		name := ctx.DefaultQuery("name", "default")

		if key, err = c.apiKey.CreateApiKey(userID, name); err != nil {
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
