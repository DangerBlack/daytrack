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

type resolvedTarget struct {
	track *models.Track
}

func (s *Service) resolveTrack(userID int64, username string, trackName string, requireWrite bool) (*resolvedTarget, error) {
	var err error
	var user, user2 *models.User
	var track *models.Track

	if user, err = s.db.GetUserByID(userID); err != nil {
		return nil, err
	}

	if user2, err = s.db.GetUserByName(username); err != nil {
		return nil, err
	}

	if track, err = s.db.GetTrackByUserIDAndName(user2.ID, trackName); err != nil {
		return nil, err
	}

	isOwner := user.Username == username
	canRead := isOwner || track.Visibility == models.TrackVisibilityPublicRead || track.Visibility == models.TrackVisibilityPublicWrite
	canWrite := isOwner || track.Visibility == models.TrackVisibilityPublicWrite

	if requireWrite && !canWrite {
		return nil, ErrorEventCannotBeCalledByYou
	}
	if !requireWrite && !canRead {
		return nil, ErrorEventCannotBeCalledByYou
	}

	return &resolvedTarget{track: track}, nil
}

func (s *Service) CreateEvent(userID int64, username string, trackName string, quantity int, createdAt *time.Time) error {
	target, err := s.resolveTrack(userID, username, trackName, true)
	if err != nil {
		return err
	}

	return s.db.InsertEvent(target.track.ID, quantity, createdAt)
}

func (s *Service) ListEventsByMonth(userID int64, username string, trackName string, after *time.Time) ([]models.Day, error) {
	target, err := s.resolveTrack(userID, username, trackName, false)
	if err != nil {
		return nil, err
	}

	return s.db.GetEventsByTrackIDGroupByMonth(target.track.ID, after)
}

func (s *Service) ListEventsByDays(userID int64, username string, trackName string, after *time.Time) ([]models.Day, error) {
	target, err := s.resolveTrack(userID, username, trackName, false)
	if err != nil {
		return nil, err
	}

	return s.db.GetEventsByTrackIDGroupByDay(target.track.ID, after)
}

func (s *Service) ListEvents(userID int64, username string, trackName string, after *time.Time) ([]models.Day, error) {
	target, err := s.resolveTrack(userID, username, trackName, false)
	if err != nil {
		return nil, err
	}

	return s.db.GetEventsByTrackID(target.track.ID, after)
}
