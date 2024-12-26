package postgrestorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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

func (m *PostgreStorage) CreateUser(
	ctx *context.Context,
	user *bizmodels.User,
) error {
	_, err := m.conn.Exec(
		*ctx,
		"INSERT INTO user (login, password) VALUES ($1, $2)",
		user.GetLogin(),
		user.GetPassword())
	if err != nil {
		return fmt.Errorf(
			"CreateUser->INSERT INTO error: %w", err)
	}

	return nil
}

func (m *PostgreStorage) GetUser(
	ctx *context.Context,
	login string,
) (*bizmodels.User, error) {
	user := &bizmodels.User{}

	var (
		outID          int32
		outLogin       string
		outPass        string
		outCreateddate time.Time
	)

	err := m.conn.QueryRow(
		*ctx,
		"select id, login, password, createddate"+
			"from user"+
			"where login=$1 LIMIT 1",
		login).Scan(&outID, &outLogin, &outPass, &outCreateddate)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}

	if err != nil {
		return nil,
			fmt.Errorf("GetUser->m.conn.QueryRow %w",
				err)
	}

	user.SetUser(outID,
		outPass,
		outLogin,
		outCreateddate,
	)

	return user, nil
}

func (m *PostgreStorage) CreateOrder(
	ctx *context.Context,
	order *bizmodels.Order,
) error {
	_, err := m.conn.Exec(
		*ctx,
		"INSERT INTO order (identifier,client) VALUES ($1,$2)",
		order.GetIdentifier(), order.GetClient().GetID())
	if err != nil {
		return fmt.Errorf(
			"CreateOrder->INSERT INTO error: %w", err)
	}

	return nil
}

func (m *PostgreStorage) GetOrder(
	ctx *context.Context,
	ident string,
) (*bizmodels.Order, error) {
	order := &bizmodels.Order{}
	user := &bizmodels.User{}

	var (
		outOrderID          int32
		outOrderIdentifier  string
		outOrderCreateddate time.Time

		outUserID          int32
		outUserLogin       string
		outUserPass        string
		outUserCreateddate time.Time
	)

	err := m.conn.QueryRow(
		*ctx,
		"select o.id,o.identifier,o.createddate,"+
			" u.id, u.login, u.password, u.createddate"+
			" from order o"+
			" left join user u on u.id = order.client"+
			" where o.identifier=$1 LIMIT 1",
		ident).Scan(&outOrderID,
		&outOrderIdentifier,
		&outOrderCreateddate,
		&outUserID,
		&outUserLogin,
		&outUserPass,
		&outUserCreateddate)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}

	if err != nil {
		return nil,
			fmt.Errorf("GetOrder->m.conn.QueryRow %w",
				err)
	}

	user.SetUser(outUserID,
		outUserPass,
		outUserLogin,
		outUserCreateddate)

	order.SetOrder(
		outOrderID,
		outOrderIdentifier,
		user,
		outOrderCreateddate)

	return order, nil
}
