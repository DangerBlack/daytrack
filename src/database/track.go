package database

import "512b.it/daytrack/src/models"

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

func (d *Database) GetTrackByID(id int64) (*models.Track, error) {
	track := models.Track{
		ID: id,
	}

	err := d.db.QueryRow(`
		SELECT
			user_id,
			name,
			description,
			visibility,
			status,
			created_at,
			delete_at
		FROM tracks
		WHERE id = ?;
	`, id).Scan(&track.UserID, &track.Name, &track.Description, &track.Visibility, &track.Status, &track.CreatedAt, &track.DeleteAt)

	if err != nil {
		return nil, err
	}

	return &track, nil
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

	return tracks, nil
}
