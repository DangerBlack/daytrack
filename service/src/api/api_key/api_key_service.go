package api_key

import (
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
	logger := utils.InitServiceLogger("ApiKeyService")

	return &Service{
		db:            db,
		configuration: configuration,
		logger:        logger,
	}
}

func (s *Service) CreateApiKey(userID int64, name string) (*string, error) {
	return s.db.InsertAPIKey(userID, name)
}

func (s *Service) ListApiKeys(userID int64) ([]models.ApiKey, error) {
	var err error
	var keys []models.ApiKey

	if keys, err = s.db.ListAPIKeys(userID); err != nil {
		return nil, err
	}

	for i, key := range keys {
		keys[i].Key = key.Key[:3] + "..." + key.Key[len(key.Key)-3:]
	}

	return keys, nil
}

func (s *Service) DeleteApiKey(userID int64, keyID int64) error {
	return s.db.DeleteAPIKey(userID, keyID)
}
