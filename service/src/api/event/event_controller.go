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
		v1.POST("/events/:username/:track_name", utils.RateLimit(3, time.Second), c.createEventRoute())
		v1.GET("/events/:username/:track_name", c.listEventRoute())
	}
}

func (c *EventController) injectAuthenticatedRoutes() {
}

// @Tags event
// @Security ApiKeyAuth
// @Schemes https
// @Router /v1/events/{username}/{track_name} [POST]
// @Summary Create an event
// @Description Create an event for the authenticated user
// @Param username path string true "Username"
// @Param track_name path string true "Track name"
// @Param created_at query string false "Created at"
// @Param quantity query int false "Quantity"
// @Accept json
// @Produce json
// @Success 201 {object} models.ResponseModel
func (c *EventController) createEventRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var userID int64
		var createdAt *time.Time

		username := ctx.Param("username")
		trackName := ctx.Param("track_name")
		createdAtString := ctx.DefaultQuery("created_at", "")
		quantity, err := strconv.Atoi(ctx.DefaultQuery("quantity", "1"))

		if err != nil {
			ctx.JSON(400, models.NewError(models.ErrorBadRequest, "invalid quantity"))
			return
		}

		if createdAtString != "" {
			createdAtx, err := time.Parse(time.RFC3339, createdAtString)
			if err != nil {
				c.logger(ctx).Err(err).Msg("Failed to parse created_at")
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

			if err == ErrorEventCannotBeCalledByYou {
				ctx.JSON(403, models.NewError(models.ErrorForbidden, "event cannot be called by you"))
				return
			}

			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to create event"))
			return
		}

		ctx.JSON(201, models.NewSuccess("event registered", ""))
	}
}

// @Tags event
// @Security ApiKeyAuth
// @Schemes https
// @Router /v1/events/{username}/{track_name} [GET]
// @Summary List events
// @Description List all events for a track
// @Param username path string true "Username"
// @Param track_name path string true "Track name"
// @Param created_at query string false "Created at"
// @Param quantity query int false "Quantity"
// @Accept json
// @Produce json
// @Success 200 {object} models.List[models.Day]
func (c *EventController) listEventRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var userID int64
		var events []models.Day
		var after *time.Time

		username := ctx.Param("username")
		trackName := ctx.Param("track_name")
		afterString := utils.EmptyIsNull(ctx.DefaultQuery("after", ""))

		listBy := models.ListBy(ctx.DefaultQuery("list_by", string(models.ListByDay)))

		if userID, err = utils.GetAuthenticatedUserID(ctx); err != nil {
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to get authenticated user id"))
			return
		}

		if afterString != nil {
			afterX, err := time.Parse(time.RFC3339, *afterString)
			if err != nil {
				ctx.JSON(400, models.NewError(models.ErrorBadRequest, "invalid created_at"))
				return
			}

			after = &afterX
		}

		switch listBy {
		case models.ListByDay:
			if events, err = c.event.ListEventsByDays(userID, username, trackName, after); err != nil {
				c.logger(ctx).Err(err).Msg("Failed to list events grouped by day")
				ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to list events"))
				return
			}
		case models.ListByMonth:
			if events, err = c.event.ListEventsByMonth(userID, username, trackName, after); err != nil {
				c.logger(ctx).Err(err).Msg("Failed to list events grouped by month")
				ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to list events"))
				return
			}
		case models.ListByRaw:
			if events, err = c.event.ListEvents(userID, username, trackName, after); err != nil {
				c.logger(ctx).Err(err).Msg("Failed to list events raws")
				ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to list events"))
				return
			}
		default:
			ctx.JSON(400, models.NewError(models.ErrorBadRequest, "invalid list_by"))
			return
		}

		ctx.JSON(200, models.List[models.Day]{
			Items: events,
		})
	}
}
