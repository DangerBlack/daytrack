package api_key

import (
	"errors"

	"512b.it/daytrack/src/database"
	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
)

var ErrorOpaqueChallengeVerificationFailed error = errors.New("opaque challenge verification failed")

type Service struct {
	db            *database.Database
	configuration models.Configuration
	logger        utils.ContextLogger
}

func New(db *database.Database, configuration models.Configuration) *Service {
	logger := utils.InitServiceLogger("ApiKeyService")

	return &Service{
		db:            db,
		configuration: configuration,
		logger:        logger,
	}
}

func (s *Service) CreateApiKey(userID int, name string) (*string, error) {
	return s.db.InsertAPIKey(userID, name)
}
