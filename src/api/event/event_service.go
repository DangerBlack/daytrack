package event

import (
	"errors"
	"time"

	"512b.it/daytrack/src/database"
	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
)

var ErrorEventCannotBeCalledByYou = errors.New("this event cannot be called by you")

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
	var user, user2 *models.User
	var track *models.Track

	if user, err = s.db.GetUserByID(userID); err != nil {
		return err
	}

	if user2, err = s.db.GetUserByName(username); err != nil {
		return err
	}

	if track, err = s.db.GetTrackByUserIDAndName(user2.ID, trackName); err != nil {
		return err
	}

	if user.Username != username && track.Visibility != models.TrackVisibilityPublicWrite {
		return ErrorEventCannotBeCalledByYou
	}

	return s.db.InsertEvent(track.ID, quantity, createdAt)
}

func (s *Service) ListEventsByDays(userID int64, username string, trackName string, after *time.Time) ([]models.Day, error) {
	var err error
	var user, user2 *models.User
	var track *models.Track
	var events []models.Day

	if user, err = s.db.GetUserByID(userID); err != nil {
		return nil, err
	}

	if user2, err = s.db.GetUserByName(username); err != nil {
		return nil, err
	}

	if track, err = s.db.GetTrackByUserIDAndName(user2.ID, trackName); err != nil {
		return nil, err
	}

	if user.Username != username && (track.Visibility != models.TrackVisibilityPublicRead && track.Visibility != models.TrackVisibilityPublicWrite) {
		return nil, ErrorEventCannotBeCalledByYou
	}

	if events, err = s.db.GetEventsByTrackIDGroupByDay(track.ID, after); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Service) ListEvents(userID int64, username string, trackName string, after *time.Time) ([]models.Day, error) {
	var err error
	var user, user2 *models.User
	var track *models.Track
	var events []models.Day

	if user, err = s.db.GetUserByID(userID); err != nil {
		return nil, err
	}

	if user2, err = s.db.GetUserByName(username); err != nil {
		return nil, err
	}

	if track, err = s.db.GetTrackByUserIDAndName(user2.ID, trackName); err != nil {
		return nil, err
	}

	if user.Username != username && (track.Visibility != models.TrackVisibilityPublicRead && track.Visibility != models.TrackVisibilityPublicWrite) {
		return nil, ErrorEventCannotBeCalledByYou
	}

	if events, err = s.db.GetEventsByTrackID(track.ID, after); err != nil {
		return nil, err
	}

	return events, nil
}
