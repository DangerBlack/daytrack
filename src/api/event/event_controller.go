package event

import (
	"strconv"
	"time"

	"512b.it/daytrack/src/api/middleware"
	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
	"github.com/gin-gonic/gin"
)

type EventController struct {
	unauthenticatedRoute *gin.RouterGroup
	authenticatedRoute   *gin.RouterGroup

	configuration models.Configuration
	event         *Service
	logger        utils.ContextLogger
}

func Inject(
	unauthenticatedRoute *gin.RouterGroup,
	authenticatedRoute *gin.RouterGroup,
	event *Service,
	configuration models.Configuration,
) {
	controller := new(unauthenticatedRoute, authenticatedRoute, event, configuration)
	controller.injectUnauthenticatedRoutes()
	controller.injectAuthenticatedRoutes()
}

func new(
	unauthenticatedRoute *gin.RouterGroup,
	authenticatedRoute *gin.RouterGroup,
	event *Service,
	configuration models.Configuration,
) *EventController {
	logger := utils.InitServiceLogger("EventController")

	return &EventController{
		unauthenticatedRoute: unauthenticatedRoute,
		authenticatedRoute:   authenticatedRoute,
		event:                event,
		configuration:        configuration,
		logger:               logger,
	}
}

func (c *EventController) injectUnauthenticatedRoutes() {
	v1 := c.unauthenticatedRoute.Group("v1", middleware.AuthApiKeyGuards(c.configuration, c.event.db))
	{
		v1.POST("/events/:username/:track_name", c.createEventRoute())
		v1.GET("/events/:username/:track_name", c.listEventRoute())

	}
}

func (c *EventController) injectAuthenticatedRoutes() {
}

func (c *EventController) createEventRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var userID int64
		var createdAt *time.Time

		username := ctx.Param("username")
		trackName := ctx.Param("track_name")
		createdAtString := utils.EmptyIsNull(ctx.DefaultQuery("created_at", ""))
		quantity, err := strconv.Atoi(ctx.DefaultQuery("quantity", "1"))

		if err != nil {
			ctx.JSON(400, models.NewError(models.ErrorBadRequest, "invalid quantity"))
			return
		}

		if createdAtString != nil {
			createdAtx, err := time.Parse(time.RFC3339, *createdAtString)
			if err != nil {
				ctx.JSON(400, models.NewError(models.ErrorBadRequest, "invalid created_at"))
				return
			}

			createdAt = &createdAtx
		}

		if userID, err = utils.GetAuthenticatedUserID(ctx); err != nil {
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to get authenticated user id"))
			return
		}

		if err = c.event.CreateEvent(userID, username, trackName, quantity, createdAt); err != nil {
			c.logger(ctx).Err(err).Msg("Failed to create event")
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to create event"))
			return
		}

		ctx.JSON(201, models.NewSuccess("event registered", ""))
	}
}

func (c *EventController) listEventRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var userID int64

		username := ctx.Param("username")
		trackName := ctx.Param("track_name")

		if userID, err = utils.GetAuthenticatedUserID(ctx); err != nil {
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to get authenticated user id"))
			return
		}

		events, err := c.event.ListEvents(userID, username, trackName)
		if err != nil {
			c.logger(ctx).Err(err).Msg("Failed to list events")
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to list events"))
			return
		}

		ctx.JSON(200, models.List[models.Day]{
			Items: events,
		})
	}
}
