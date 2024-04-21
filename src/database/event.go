package database

import (
	"time"

	"512b.it/daytrack/src/models"
)

func (d *Database) InsertEvent(trackID int64, quantity int, createdAt *time.Time) error {
	_, err := d.db.Exec(`
		INSERT INTO events (track_id, quantity, created_at)
		VALUES (?, ?, ?);
	`, trackID, quantity, createdAt)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetEventsByTrackID(trackID int64) ([]models.Day, error) {
	rows, err := d.db.Query(`
		SELECT
			DATE(created_at),
			SUM(quantity)
		FROM events
		WHERE track_id = ?
		GROUP BY DATE(created_at)
		ORDER BY DATE(created_at) DESC;
	`, trackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Day
	for rows.Next() {
		var event models.Day
		err := rows.Scan(&event.Date, &event.Quantity)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}
