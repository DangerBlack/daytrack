package database

import (
	"512b.it/daytrack/src/models"
)

func (d *Database) InsertUser(username, email, password, salt string) (int, error) {
	res, err := d.db.Exec(`
		INSERT INTO users (username, email, password, salt)
		VALUES (?, ?, ?, ?);
	`, username, email, password, salt)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (d *Database) GetUserByID(id int) (*models.User, error) {
	user := models.User{
		ID: id,
	}

	err := d.db.QueryRow(`
		SELECT
			username,
			email,
			password,
			salt
		FROM users
		WHERE id = ?;
	`, id).Scan(&user.Username, &user.Email, &user.Password, &user.Salt)

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
			password,
			salt
		FROM users
		WHERE email = ?;
	`, email).Scan(&user.ID, &user.Username, &user.Password, &user.Salt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
