package database

import (
	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
)

func (d *Database) InsertAPIKey(userID int, name string) (*string, error) {
	key := utils.GenerateRandomString(32)

	_, err := d.db.Exec(`
		INSERT INTO api_keys (user_id, key, name)
		VALUES (?, ?, ?);
	`, userID, key, name)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func (d *Database) GetAPIKey(key string) (*models.ApiKey, error) {
	apiKey := models.ApiKey{
		Key: key,
	}

	err := d.db.QueryRow(`
		SELECT 
			id,
			user_id,
			created_at
		FROM api_keys
		WHERE
			key = ?
			AND delete_at IS NULL;
	`, key).Scan(&apiKey.ID, &apiKey.UserID, apiKey.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}

func (d *Database) DeleteAPIKey(key string) error {
	_, err := d.db.Exec(`
		UPDATE api_keys
		SET delete_at = CURRENT_TIMESTAMP
		WHERE key = ?;
	`, key)
	if err != nil {
		return err
	}

	return nil
}
