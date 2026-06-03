package database

import (
	"database/sql"

	"512b.it/daytrack/src/models"
)

func (d *Database) InsertTrack(userID int64, name string, description string, visibility models.TrackVisibility, status models.TrackStatus) (*int64, error) {
	res, err := d.db.Exec(`
		INSERT INTO tracks (user_id, name, description, visibility, status)
		VALUES (?, ?, ?, ?, ?);
	`, userID, name, description, visibility, status)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (d *Database) GetTrackByUserIDAndName(userID int64, trackName string) (*models.Track, error) {
	track := models.Track{
		UserID: userID,
		Name:   trackName,
	}

	err := d.db.QueryRow(`
		SELECT
			id,
			description,
			visibility,
			status,
			created_at,
			delete_at
		FROM tracks
		WHERE user_id = ? and name = ?;
	`, userID, trackName).Scan(&track.ID, &track.Description, &track.Visibility, &track.Status, &track.CreatedAt, &track.DeleteAt)

	if err != nil {
		return nil, err
	}

	return &track, nil
}

func (d *Database) ListTracks(userID int64) ([]models.Track, error) {
	rows, err := d.db.Query(`
		SELECT
			id,
			name,
			description,
			visibility,
			status,
			created_at
		FROM tracks
		WHERE
			user_id = ?
			AND delete_at IS NULL;
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []models.Track
	for rows.Next() {
		var track models.Track
		err := rows.Scan(&track.ID, &track.Name, &track.Description, &track.Visibility, &track.Status, &track.CreatedAt)
		if err != nil {
			return nil, err
		}

		tracks = append(tracks, track)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (d *Database) UpdateTrack(userID int64, trackName string, name *string, description *string, visibility *models.TrackVisibility, status *models.TrackStatus) error {
	_, err := d.db.Exec(`
		UPDATE tracks
		SET
			name =  COALESCE(:name, name),
			description = COALESCE(:description, description),
			visibility = COALESCE(:visibility, visibility),
			status = COALESCE(:status, status)
		WHERE
			user_id = :user_id
			AND name = :track_name;
	`,
		sql.Named("name", GetNullString(name)),
		sql.Named("description", GetNullString(description)),
		sql.Named("visibility", GetNullVisibility(visibility)),
		sql.Named("status", GetNullStatus(status)),
		sql.Named("user_id", userID),
		sql.Named("track_name", trackName),
	)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) DeleteTrack(userID int64, trackName string) error {
	_, err := d.db.Exec(`
		UPDATE tracks
		SET
			delete_at = CURRENT_TIMESTAMP
		WHERE
			user_id = ?
			AND name = ?;
	`, userID, trackName)
	if err != nil {
		return err
	}

	return nil
}
