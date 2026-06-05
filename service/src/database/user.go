package database

import (
	"database/sql"
	"errors"

	"512b.it/daytrack/src/models"
)

func (d *Database) InsertUser(username, email, publicKey, salt string) (int64, error) {
	res, err := d.db.Exec(`
		INSERT INTO users (username, email, public_key, salt)
		VALUES (?, ?, ?, ?);
	`, username, email, publicKey, salt)
	if err != nil {
		return 0, ParseError(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *Database) GetUserByID(id int64) (*models.User, error) {
	user := models.User{
		ID: id,
	}

	err := d.db.QueryRow(`
		SELECT
			username,
			email,
			public_key,
			salt
		FROM users
		WHERE id = ?;
	`, id).Scan(&user.Username, &user.Email, &user.PublicKey, &user.Salt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrorNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (d *Database) GetUserByName(username string) (*models.User, error) {
	user := models.User{
		Username: username,
	}

	err := d.db.QueryRow(`
		SELECT
			id,
			email,
			public_key,
			salt
		FROM users
		WHERE username = ?;
	`, username).Scan(&user.ID, &user.Email, &user.PublicKey, &user.Salt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrorNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (d *Database) GetUserByEmail(email string) (*models.User, error) {
	user := models.User{
		Email: email,
	}

	err := d.db.QueryRow(`
		SELECT
			id,
			username,
			public_key,
			salt
		FROM users
		WHERE email = ?;
	`, email).Scan(&user.ID, &user.Username, &user.PublicKey, &user.Salt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrorNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
