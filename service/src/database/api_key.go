package database

import (
	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
)

func (d *Database) InsertAPIKey(userID int64, name string) (int64, string, error) {
	key := utils.GenerateRandomString(32)

	res, err := d.db.Exec(`
		INSERT INTO api_keys (user_id, key, name)
		VALUES (?, ?, ?);
	`, userID, key, name)
	if err != nil {
		return 0, "", err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, "", err
	}

	return id, key, nil
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
	`, key).Scan(&apiKey.ID, &apiKey.UserID, &apiKey.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}

func (d *Database) ListAPIKeys(userID int64) ([]models.ApiKey, error) {
	rows, err := d.db.Query(`
		SELECT 
			id,
			user_id,
			name,
			key,
			created_at
		FROM api_keys
		WHERE
			user_id = ?
			AND delete_at IS NULL;
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apiKeys := make([]models.ApiKey, 0)
	for rows.Next() {
		apiKey := models.ApiKey{}

		if err := rows.Scan(&apiKey.ID, &apiKey.UserID, &apiKey.Name, &apiKey.Key, &apiKey.CreatedAt); err != nil {
			return nil, err
		}

		apiKeys = append(apiKeys, apiKey)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return apiKeys, nil
}

func (d *Database) DeleteAPIKey(userID int64, keyID int64) error {
	_, err := d.db.Exec(`
		UPDATE api_keys
		SET delete_at = CURRENT_TIMESTAMP
		WHERE 
			user_id = ?
		AND
			id = ?;
	`, userID, keyID)
	if err != nil {
		return err
	}

	return nil
}
