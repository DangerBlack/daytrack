package database

import (
	"database/sql"
	"errors"
	"strings"

	"512b.it/daytrack/src/models"
	_ "github.com/mattn/go-sqlite3"
)

var ErrorDuplicate = errors.New("duplicate")
var ErrorNotFound = errors.New("not found")

type Database struct {
	db *sql.DB
}

func NewDatabase(dbPath string) *Database {
	if dbPath == "" {
		dbPath = "./archive/database.db"
	}

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(1)

	database := &Database{
		db: db,
	}
	database.InitDatabase()

	return database
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) InitDatabase() {
	_, err := d.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			email TEXT NOT NULL,
			public_key TEXT NOT NULL,
			salt TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(username),
			UNIQUE(email)
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
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
			status TEXT NOT NULL DEFAULT 'enabled' CHECK (status IN ('enabled', 'disabled')),
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			quantity INTEGER NOT NULL,
			FOREIGN KEY(track_id) REFERENCES tracks(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		panic(err)
	}

}

func GetNullString(s *string) sql.NullString {
	nullString := sql.NullString{
		Valid: s != nil,
	}
	if s != nil {
		nullString.String = *s
	}

	return nullString
}

func GetNullStatus(s *models.TrackStatus) sql.NullString {
	nullString := sql.NullString{
		Valid: s != nil,
	}
	if s != nil {
		nullString.String = string(*s)
	}

	return nullString
}

func GetNullVisibility(s *models.TrackVisibility) sql.NullString {
	nullString := sql.NullString{
		Valid: s != nil,
	}
	if s != nil {
		nullString.String = string(*s)
	}

	return nullString
}

func ParseError(err error) error {
	if err == nil {
		return nil
	}

	if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
		return ErrorDuplicate
	}

	return err
}
