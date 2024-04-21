package track

import (
	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
	"github.com/gin-gonic/gin"
)

type TrackController struct {
	unauthenticatedRoute *gin.RouterGroup
	authenticatedRoute   *gin.RouterGroup

	configuration models.Configuration
	track         *Service
	logger        utils.ContextLogger
}

func Inject(
	unauthenticatedRoute *gin.RouterGroup,
	authenticatedRoute *gin.RouterGroup,
	track *Service,
	configuration models.Configuration,
) {
	controller := new(unauthenticatedRoute, authenticatedRoute, track, configuration)
	controller.injectUnauthenticatedRoutes()
	controller.injectAuthenticatedRoutes()
}

func new(
	unauthenticatedRoute *gin.RouterGroup,
	authenticatedRoute *gin.RouterGroup,
	track *Service,
	configuration models.Configuration,
) *TrackController {
	logger := utils.InitServiceLogger("TrackController")

	return &TrackController{
		unauthenticatedRoute: unauthenticatedRoute,
		authenticatedRoute:   authenticatedRoute,
		track:                track,
		configuration:        configuration,
		logger:               logger,
	}
}

func (c *TrackController) injectUnauthenticatedRoutes() {
}

func (c *TrackController) injectAuthenticatedRoutes() {
	v1 := c.authenticatedRoute.Group("v1")
	{
		v1.POST("/tracks", c.createTrackRoute())
		v1.GET("/tracks", c.listTracksRoute())
		v1.PATCH("/tracks/:track_name", c.updateTracksRoute())
	}
}

// @Tags track
// @Security TokenAuth
// @Schemes https
// @Router /v1/tracks [POST]
// @Summary Create a track
// @Description Create a track
// @Accept json
// @Param request body models.Track true "Track object"
// @Produce json
// @Success 201 {object} models.ResponseModel
func (c *TrackController) createTrackRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var userID int64

		if userID, err = utils.GetAuthenticatedUserID(ctx); err != nil {
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to get authenticated user id"))
			return
		}

		var track models.Track

		if err := ctx.ShouldBindJSON(&track); err != nil {
			c.logger(ctx).Err(err).Msg("Invalid request")
			ctx.JSON(400, models.NewError(models.ErrorBadRequest, "invalid request"))
			return
		}

		c.logger(ctx).Debug().Msgf("track: %v", track)

		if track.Visibility == "" {
			track.Visibility = models.TrackVisibilityPrivate
		}

		if track.Status == "" {
			track.Status = models.TrackStatusEnabled
		}

		if _, err = c.track.CreateTrack(ctx, int64(userID), track); err != nil {
			c.logger(ctx).Err(err).Msg("Error while creating a track")
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to create track"))
			return
		}

		ctx.JSON(201, models.NewSuccess("track created", ""))
	}
}

// @Tags track
// @Security TokenAuth
// @Schemes https
// @Router /v1/tracks [GET]
// @Summary List tracks
// @Description List all tracks for the authenticated user
// @Accept json
// @Produce json
// @Success 201 {object} models.List[models.Track]
func (c *TrackController) listTracksRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var userID int64
		var tracks []models.Track

		if userID, err = utils.GetAuthenticatedUserID(ctx); err != nil {
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to get authenticated user id"))
			return
		}

		if tracks, err = c.track.ListTracks(ctx, userID); err != nil {
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to list tracks"))
			return
		}

		ctx.JSON(200, models.List[models.Track]{
			Items: tracks,
		})
	}
}

func (c *TrackController) updateTracksRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var userID int64

		if userID, err = utils.GetAuthenticatedUserID(ctx); err != nil {
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to get authenticated user id"))
			return
		}

		trackName := ctx.Param("track_name")

		var track models.Track

		if err := ctx.ShouldBindJSON(&track); err != nil {
			c.logger(ctx).Err(err).Msg("Invalid request")
			ctx.JSON(400, models.NewError(models.ErrorBadRequest, "invalid request"))
			return
		}

		if err = c.track.UpdateTrack(ctx, userID, trackName, track); err != nil {
			c.logger(ctx).Err(err).Msg("Error while updating a track")
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to update track"))
			return
		}

		ctx.JSON(200, models.NewSuccess("track updated", ""))
	}
}
