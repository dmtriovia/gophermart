package userstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
	"github.com/jackc/pgx/v5"
)

type UserStorage struct {
	conn *pgx.Conn
}

func (m *UserStorage) Initiate(
	conn *pgx.Conn,
) {
	m.conn = conn
}

const defUserData = "u.id,u.login,u.password,u.createddate"

func (m *UserStorage) CreateUser(
	ctx *context.Context,
	user *usermodel.User,
) error {
	var lastInsertID int32

	err := m.conn.QueryRow(
		*ctx,
		"INSERT INTO users (login, password)"+
			" VALUES ($1, $2) RETURNING id",
		user.GetLogin(),
		user.GetPassword()).Scan(&lastInsertID)
	if err != nil {
		return fmt.Errorf(
			"CreateUser->Scan: %w", err)
	}

	user.SetID(lastInsertID)

	return nil
}

func (m *UserStorage) GetUser(
	ctx *context.Context,
	login string,
) (*usermodel.User, error) {
	user := &usermodel.User{}

	var (
		outID             int32
		outLogin, outPass string
		outCreateddate    time.Time
	)

	err := m.conn.QueryRow(
		*ctx,
		"select "+defUserData+
			" from users u"+
			" where login=$1 LIMIT 1",
		login).Scan(&outID, &outLogin, &outPass,
		&outCreateddate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}

		return nil,
			fmt.Errorf("GetUser->m.conn.QueryRow %w",
				err)
	}

	user.SetUser(outID,
		outLogin,
		outPass,
		outCreateddate,
	)

	return user, nil
}
