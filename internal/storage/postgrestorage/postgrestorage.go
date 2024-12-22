package postgrestorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels"
	"github.com/jackc/pgx/v5"
)

type PostgreStorage struct {
	conn *pgx.Conn
}

func (m *PostgreStorage) Initiate(
	conn *pgx.Conn,
) {
	m.conn = conn
}

func (m *PostgreStorage) GetUser(
	ctx *context.Context,
	login string,
) (*bizmodels.User, error) {
	user := &bizmodels.User{}

	var (
		outLogin string
		outPass  string
	)

	err := m.conn.QueryRow(
		*ctx,
		"select login, password from user where login=$1 LIMIT 1",
		login).Scan(&outLogin, &outPass)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}

	if err != nil {
		return nil,
			fmt.Errorf("GetUser->m.conn.QueryRow %w",
				err)
	}

	user.Password = outPass
	user.Login = outLogin

	return user, nil
}

func (m *PostgreStorage) CreateUser(
	ctx *context.Context,
	user *bizmodels.User,
) error {
	_, err := m.conn.Exec(
		*ctx,
		"INSERT INTO user (login, password) VALUES ($1, $2)",
		user.Login,
		user.Password)
	if err != nil {
		return fmt.Errorf(
			"CreateUser->INSERT INTO error: %w", err)
	}

	return nil
}
