package track

import (
	"512b.it/daytrack/src/database"
	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
	"github.com/gin-gonic/gin"
)

type Service struct {
	db            *database.Database
	configuration models.Configuration
	logger        utils.ContextLogger
}

func New(db *database.Database, configuration models.Configuration) *Service {
	logger := utils.InitServiceLogger("TrackService")

	return &Service{
		db:            db,
		configuration: configuration,
		logger:        logger,
	}
}

func (s *Service) CreateTrack(ctx *gin.Context, userID int64, track models.CreateTrack) (*int64, error) {
	return s.db.InsertTrack(userID, track.Name, track.Description, track.Visibility, track.Status)
}

func (s *Service) ListTracks(ctx *gin.Context, userID int64) ([]models.Track, error) {
	return s.db.ListTracks(userID)
}

func (s *Service) UpdateTrack(ctx *gin.Context, userID int64, trackName string, track models.Track) error {
	var visibility models.TrackVisibility
	var status models.TrackStatus

	if track.Visibility != "" {
		visibility = track.Visibility
	}

	if track.Status != "" {
		status = track.Status
	}

	return s.db.UpdateTrack(userID, trackName, utils.EmptyIsNull(track.Name), utils.EmptyIsNull(track.Description), &visibility, &status)
}
