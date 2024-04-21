package event

import (
	"errors"
	"time"

	"512b.it/daytrack/src/database"
	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
)

type Service struct {
	db            *database.Database
	configuration models.Configuration
	logger        utils.ContextLogger
}

func New(db *database.Database, configuration models.Configuration) *Service {
	logger := utils.InitServiceLogger("EventService")

	return &Service{
		db:            db,
		configuration: configuration,
		logger:        logger,
	}
}

func (s *Service) CreateEvent(userID int64, username string, trackName string, quantity int, createdAt *time.Time) error {
	var err error
	var user *models.User
	var track *models.Track

	if user, err = s.db.GetUserByID(userID); err != nil {
		return err
	}

	if user.Username != username {
		return errors.New("this event cannot be called by you")
	}

	if track, err = s.db.GetTrackByUserIDAndName(user.ID, trackName); err != nil {
		return err
	}

	return s.db.InsertEvent(track.ID, quantity, createdAt)
}

func (s *Service) ListEvents(userID int64, username string, trackName string) ([]models.Day, error) {
	var err error
	var user *models.User
	var track *models.Track
	var events []models.Day

	if user, err = s.db.GetUserByID(userID); err != nil {
		return nil, err
	}

	if user.Username != username {
		return nil, errors.New("this event cannot be called by you")
	}

	if track, err = s.db.GetTrackByUserIDAndName(user.ID, trackName); err != nil {
		return nil, err
	}

	if events, err = s.db.GetEventsByTrackID(track.ID); err != nil {
		return nil, err
	}

	return events, nil
}
