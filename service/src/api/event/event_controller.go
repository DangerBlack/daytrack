package event

import (
	"errors"
	"strconv"
	"time"

	"512b.it/daytrack/src/api/middleware"
	"512b.it/daytrack/src/database"
	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
	"github.com/gin-gonic/gin"
)

type EventController struct {
	unauthenticatedRoute *gin.RouterGroup
	authenticatedRoute   *gin.RouterGroup

	db            *database.Database
	configuration models.Configuration
	event         *Service
	logger        utils.ContextLogger
}

func Inject(
	unauthenticatedRoute *gin.RouterGroup,
	authenticatedRoute *gin.RouterGroup,
	db *database.Database,
	event *Service,
	configuration models.Configuration,
) {
	controller := new(unauthenticatedRoute, authenticatedRoute, db, event, configuration)
	controller.injectUnauthenticatedRoutes()
	controller.injectAuthenticatedRoutes()
}

func new(
	unauthenticatedRoute *gin.RouterGroup,
	authenticatedRoute *gin.RouterGroup,
	db *database.Database,
	event *Service,
	configuration models.Configuration,
) *EventController {
	logger := utils.InitServiceLogger("EventController")

	return &EventController{
		unauthenticatedRoute: unauthenticatedRoute,
		authenticatedRoute:   authenticatedRoute,
		db:                   db,
		event:                event,
		configuration:        configuration,
		logger:               logger,
	}
}

func (c *EventController) injectUnauthenticatedRoutes() {
	v1 := c.unauthenticatedRoute.Group("v1")
	v1.Use(middleware.AuthApiKeyGuards(c.configuration, c.db))
	{
		v1.POST("/events/:username/:track_name", utils.RateLimit(c.configuration.RateLimitBurst, c.configuration.RateLimitInterval), c.createEventRoute())
		v1.GET("/events/:username/:track_name", c.listEventRoute())
	}
}

func (c *EventController) injectAuthenticatedRoutes() {}

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
				ctx.JSON(400, models.NewError(models.ErrorBadRequest, "invalid created_at parameter, expected RFC 3339 format"))
				return
			}
			createdAt = &createdAtx
		}

		var uid *int64
		if rawID, authErr := utils.GetAuthenticatedUserID(ctx); authErr == nil {
			uid = &rawID
		}

		if err = c.event.CreateEvent(uid, username, trackName, quantity, createdAt); err != nil {
			c.logger(ctx).Err(err).Msg("Failed to create event")
			if errors.Is(err, ErrorEventCannotBeCalledByYou) || errors.Is(err, database.ErrorNotFound) {
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
		var events []models.Day
		var after *time.Time

		username := ctx.Param("username")
		trackName := ctx.Param("track_name")
		afterString := utils.EmptyIsNull(ctx.DefaultQuery("after", ""))

		listBy := models.ListBy(ctx.DefaultQuery("list_by", string(models.ListByDay)))

		var uid *int64
		if rawID, authErr := utils.GetAuthenticatedUserID(ctx); authErr == nil {
			uid = &rawID
		}

		if afterString != nil {
			afterX, err := time.Parse(time.RFC3339, *afterString)
			if err != nil {
				ctx.JSON(400, models.NewError(models.ErrorBadRequest, "invalid after parameter, expected RFC 3339 format"))
				return
			}
			after = &afterX
		}

		switch listBy {
		case models.ListByDay:
			events, err = c.event.ListEventsByDays(uid, username, trackName, after)
		case models.ListByMonth:
			events, err = c.event.ListEventsByMonth(uid, username, trackName, after)
		case models.ListByRaw:
			events, err = c.event.ListEvents(uid, username, trackName, after)
		default:
			ctx.JSON(400, models.NewError(models.ErrorBadRequest, "invalid list_by"))
			return
		}

		if err != nil {
			c.logger(ctx).Err(err).Msg("Failed to list events")
			if errors.Is(err, ErrorEventCannotBeCalledByYou) {
				ctx.JSON(403, models.NewError(models.ErrorForbidden, "event cannot be called by you"))
			} else {
				ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to list events"))
			}
			return
		}

		ctx.JSON(200, models.List[models.Day]{
			Items: events,
		})
	}
}
