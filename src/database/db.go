package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() *Database {
	db, err := sql.Open("sqlite3", "./archive/database.db")
	if err != nil {
		panic(err)
	}

	database := &Database{
		db: db,
	}
	database.InitDatabase()

	return database
	// defer db.Close()
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) InitDatabase() {
	_, err := d.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			email TEXT NOT NULL,
			public_key TEXT NOT NULL,
			salt TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		panic(err)
	}

	_, err = d.db.Exec(`
		CREATE TABLE IF NOT EXISTS api_keys (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			key TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			delete_at DATETIME,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(key)
		);
	`)
	if err != nil {
		panic(err)
	}

	_, err = d.db.Exec(`
		CREATE TABLE IF NOT EXISTS tracks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			visibility TEXT NOT NULL DEFAULT 'private' CHECK (visibility IN ('private', 'public_r', 'public_rw')),
			status TEXT NOT NULL DEFAULT 'enabled' CHECK (visibility IN ('enabled', 'disabled')),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			delete_at DATETIME,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(user_id, name)
		);
	`)

	if err != nil {
		panic(err)
	}

	_, err = d.db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			track_id INTEGER NOT NULL,
			create_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			quantity INTEGER NOT NULL,
			FOREIGN KEY(track_id) REFERENCES tracks(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		panic(err)
	}

	_, err = d.db.Exec(`
		CREATE TABLE IF NOT EXISTS actions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			track_id INTEGER NOT NULL,
			create_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			action_type TEXT NOT NULL,
			action_body TEXT,
			condition TEXT,
			should_notify BOOLEAN DEFAULT FALSE,
			FOREIGN KEY(track_id) REFERENCES tracks(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		panic(err)
	}

	_, err = d.db.Exec(`
		CREATE TABLE IF NOT EXISTS history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			event_id INTEGER NOT NULL,
			action_id INTEGER NOT NULL,
			api_key_id INTEGER NOT NULL,
			create_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(event_id) REFERENCES events(id) ON DELETE CASCADE,
			FOREIGN KEY(action_id) REFERENCES actions(id) ON DELETE CASCADE,
			FOREIGN KEY(api_key_id) REFERENCES api_keys(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		panic(err)
	}
}
