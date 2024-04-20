package database

func (d *Database) InsertEvent(trackID int64, quantity int) error {
	_, err := d.db.Exec(`
		INSERT INTO events (track_id, quantity)
		VALUES (?, ?);
	`, trackID, quantity)
	if err != nil {
		return err
	}

	return nil
}
