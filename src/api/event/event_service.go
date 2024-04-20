package event

import (
	"errors"

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

func (s *Service) CreateEvent(userID int64, username string, trackName string, quantity int) error {
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

	return s.db.InsertEvent(track.ID, quantity)
}
