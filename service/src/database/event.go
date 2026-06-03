package database

import (
	"database/sql"
	"time"

	"512b.it/daytrack/src/models"
)

func (d *Database) InsertEvent(trackID int64, quantity int, createdAt *time.Time) error {
	createAtDB := sql.NullTime{
		Valid: createdAt != nil,
	}
	if createdAt != nil {
		createAtDB.Time = *createdAt
	}

	_, err := d.db.Exec(`
		INSERT INTO events (track_id, quantity, created_at)
		VALUES (?, ?, COALESCE(?, current_timestamp));
	`, trackID, quantity, createAtDB)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetEventsByTrackIDGroupByMonth(trackID int64, after *time.Time) ([]models.Day, error) {
	rows, err := d.db.Query(`
		SELECT
			datetime(strftime('%Y-%m-01 00:00:00', created_at)),
			SUM(quantity)
		FROM events
		WHERE 
				track_id = :track_id
			AND 
				(:after is NULL OR created_at > :after)
		GROUP BY strftime('%Y-%m', created_at)
		ORDER BY 1 DESC;
	`, sql.Named("track_id", trackID), sql.Named("after", after))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Day
	var timeDB string
	for rows.Next() {
		var event models.Day
		err := rows.Scan(&timeDB, &event.Quantity)
		if err != nil {
			return nil, err
		}

		if event.Date, err = time.Parse("2006-01-02 15:04:05", timeDB); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (d *Database) GetEventsByTrackIDGroupByDay(trackID int64, after *time.Time) ([]models.Day, error) {
	rows, err := d.db.Query(`
		SELECT
			datetime(strftime('%Y-%m-%d 00:00:00', created_at)),
			SUM(quantity)
		FROM events
		WHERE 
				track_id = :track_id
			AND 
				(:after is NULL OR created_at > :after)
		GROUP BY DATE(created_at)
		ORDER BY created_at DESC;
	`, sql.Named("track_id", trackID), sql.Named("after", after))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Day
	var timeDB string
	for rows.Next() {
		var event models.Day
		err := rows.Scan(&timeDB, &event.Quantity)
		if err != nil {
			return nil, err
		}

		if event.Date, err = time.Parse("2006-01-02 15:04:05", timeDB); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (d *Database) GetEventsByTrackID(trackID int64, after *time.Time) ([]models.Day, error) {
	rows, err := d.db.Query(`
		SELECT
			created_at,
			quantity
		FROM events
		WHERE 
				track_id = :track_id
			AND 
				(:after is NULL OR created_at > :after)
		ORDER BY created_at DESC;
	`, sql.Named("track_id", trackID), sql.Named("after", after))
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
