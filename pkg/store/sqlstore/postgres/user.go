package postgres

import (
	"database/sql"
	"errors"

	"github.com/bensmile/go-api-tdd/pkg/common"
	"github.com/bensmile/go-api-tdd/pkg/domain"
)

var (
	sqlCreateUser        = `INSERT INTO users (name, email, password) VALUES($1, $2, $3) RETURNING *`
	sqlDeleteUserById    = `DELETE FROM users WHERE id = $1`
	sqlSelectUserByEmail = `SELECT * FROM users WHERE email = $1`
	sqlSelectUserById    = `SELECT * FROM users WHERE id = $1`
	sqlDeleteAllUserss   = `DELETE FROM users`
)

func (q *postgresStore) CreateUser(user *domain.User) (*domain.User, error) {

	user.Password, _ = common.HashPassword(user.Password)

	err := q.db.QueryRow(sqlCreateUser, user.Name, user.Email, user.Password).
		Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (q *postgresStore) DeleteUserById(id int64) error {
	_, err := q.db.Exec(sqlDeleteUserById, id)
	if err != nil {
		return err
	}
	return nil
}

func (q *postgresStore) FindUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := q.db.QueryRow(sqlSelectUserByEmail, email).
		Scan(&user.Id, &user.Name, &user.Email, &user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (q *postgresStore) FindUserById(id int64) (*domain.User, error) {
	var user domain.User
	err := q.db.QueryRow(sqlSelectUserById, id).
		Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (q *postgresStore) DeleteAllUsers() error {
	_, err := q.db.Exec(sqlDeleteAllUserss)
	if err != nil {
		return err
	}
	return nil
}
