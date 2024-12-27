package postgrestorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels/ordermodel"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
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
	user *usermodel.User,
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
) (*usermodel.User, error) {
	user := &usermodel.User{}

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
	order *ordermodel.Order,
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
) (*ordermodel.Order, error) {
	order := &ordermodel.Order{}
	user := &usermodel.User{}

	var (
		outOrderID          int32
		outOrderStatus      string
		outOrderIdentifier  string
		outOrderCreateddate time.Time
		outOrderAccrual     int32

		outUserID          int32
		outUserLogin       string
		outUserPass        string
		outUserCreateddate time.Time
	)

	err := m.conn.QueryRow(
		*ctx,
		"select o.id,o.identifier,o.createddate,o.status,"+
			"o.accrual,"+
			" u.id, u.login, u.password, u.createddate"+
			" from order o"+
			" left join user u on u.id = order.client"+
			" where o.identifier=$1 LIMIT 1",
		ident).Scan(&outOrderID,
		&outOrderIdentifier,
		&outOrderCreateddate,
		&outOrderStatus,
		&outOrderAccrual,
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
		outOrderCreateddate,
		outOrderStatus, outOrderAccrual)

	return order, nil
}

func (m *PostgreStorage) GetOrdersByClient(
	ctx *context.Context,
	clientID int32,
) (*[]ordermodel.Order, *[]error, error) {
	var (
		outOrderID, outUserID, outOrderAccrual  int32
		outOrderStatus, outOrderIdentifier      string
		outUserLogin, outUserPass               string
		outUserCreateddate, outOrderCreateddate time.Time
	)

	rows, err := m.conn.Query(
		*ctx,
		"select o.id,o.identifier,o.createddate,o.status,"+
			"o.accrual,"+
			" u.id, u.login, u.password, u.createddate"+
			" from order o"+
			" left join user u on u.id = o.client"+
			" where o.client=$1"+
			" order by o.createddate desc",
		clientID)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"GetOrderByClient->m.conn.Query %w", err)
	}

	defer rows.Close()

	cnt := 0
	for rows.Next() {
		cnt++
	}

	orders := make([]ordermodel.Order, 0, cnt)
	errors := make([]error, 0, cnt)

	if cnt == 0 {
		return &orders, &errors, nil
	}

	for rows.Next() {
		order := &ordermodel.Order{}
		user := &usermodel.User{}
		err = rows.Scan(&outOrderID, &outOrderIdentifier,
			&outOrderCreateddate, &outOrderStatus, &outUserID,
			&outUserLogin, &outUserPass, &outUserCreateddate)

		if err != nil {
			errors = append(errors, err)
		} else {
			user.SetUser(outUserID, outUserPass,
				outUserLogin, outUserCreateddate)
			order.SetOrder(
				outOrderID, outOrderIdentifier, user,
				outOrderCreateddate, outOrderStatus, outOrderAccrual)

			orders = append(orders, *order)
		}
	}

	return &orders, &errors, nil
}
