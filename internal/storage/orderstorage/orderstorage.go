package orderstorage

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

type OrderStorage struct {
	conn *pgx.Conn
}

func (m *OrderStorage) Initiate(
	conn *pgx.Conn,
) {
	m.conn = conn
}

const defUserData = "u.id,u.login,u.password,u.createddate"

const defOrderData = "o.id, o.identifier, o.createddate, " +
	"o.status, o.accrual, o.points_write_off"

func (m *OrderStorage) CreateOrder(
	ctx *context.Context,
	order *ordermodel.Order,
) error {
	var lastInsertID *int32

	err := m.conn.QueryRow(
		*ctx,
		"INSERT INTO orders (identifier,client,"+
			" accrual,status) VALUES ($1,$2,$3,$4) RETURNING id",
		order.GetIdentifier(), order.GetClient().GetID(),
		0, order.GetStatus()).Scan(&lastInsertID)
	if err != nil {
		return fmt.Errorf(
			"CreateOrder->Scan: %w", err)
	}

	order.SetID(*lastInsertID)

	return nil
}

func (m *OrderStorage) UpdateStatusByID(
	ctx *context.Context,
	orderID int32,
	status string,
) (bool, error) {
	rows, err := m.conn.Exec(
		*ctx,
		"UPDATE orders SET status=$1 where id=$2",
		status,
		orderID)
	if err != nil {
		return false, fmt.Errorf(
			"UpdateStatusByID->m.conn.Exec( %w", err)
	}

	if rows.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}

func (m *OrderStorage) UpdateStatusAccrualByID(
	ctx *context.Context,
	tranz pgx.Tx,
	orderID int32,
	accrual float32,
	status string,
) (bool, error) {
	rows, err := tranz.Exec(
		*ctx,
		"UPDATE orders SET status=$1, accrual=$2 where id=$3",
		status,
		accrual,
		orderID)
	if err != nil {
		return false, fmt.Errorf(
			"UpdateStatusAccrualByID->tranz.Exec( %w", err)
	}

	if rows.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}

func (m *OrderStorage) ChangePointsWriteOffByID(
	ctx *context.Context,
	tranz pgx.Tx,
	orderID int32,
	newValuePointsWriteOff float32,
	sign string,
) (bool, error) {
	t := "points_write_off"

	rows, err := tranz.Exec(
		*ctx,
		"UPDATE orders SET "+t+"="+t+sign+"$1 where id=$2",
		newValuePointsWriteOff,
		orderID)
	if err != nil {
		return false, fmt.Errorf(
			"ChangePointsWriteOffByID->tranz.Exec( %w", err)
	}

	if rows.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}

func (m *OrderStorage) GetOrder(
	ctx *context.Context,
	ident string,
) (*ordermodel.Order, error) {
	order := &ordermodel.Order{}
	user := &usermodel.User{}

	var (
		outOrderID, outUserID                   *int32
		outOrderStatus, outOrderIdentifier      *string
		outOrderCreateddate, outUserCreateddate *time.Time
		outOrderPointsWriteOff, outOrderAccrual *float32
		outUserLogin, outUserPass               *string
	)

	err := m.conn.QueryRow(
		*ctx, "select "+defOrderData+","+defUserData+
			" from orders o"+
			" left join users u on u.id = o.client"+
			" where o.identifier=$1 LIMIT 1",
		ident).Scan(&outOrderID,
		&outOrderIdentifier,
		&outOrderCreateddate,
		&outOrderStatus,
		&outOrderAccrual,
		&outOrderPointsWriteOff,
		&outUserID,
		&outUserLogin,
		&outUserPass,
		&outUserCreateddate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}

		return nil,
			fmt.Errorf("GetOrder->m.conn.QueryRow %w",
				err)
	}

	user.SetUser(*outUserID,
		outUserLogin,
		outUserPass,
		outUserCreateddate)

	order.SetOrder(
		*outOrderID,
		outOrderIdentifier,
		user,
		outOrderCreateddate,
		outOrderStatus, outOrderAccrual,
		outOrderPointsWriteOff)

	return order, nil
}

func (m *OrderStorage) GetOrdersByClient(
	ctx *context.Context,
	clientID int32,
) (*[]ordermodel.Order, *[]error, error) {
	var (
		outOrderID, outUserID                   *int32
		outOrderStatus, outOrderIdentifier      *string
		outUserLogin, outUserPass               *string
		outUserCreateddate, outOrderCreateddate *time.Time
		outOrderPointsWriteOff, outOrderAccrual *float32
	)

	rows, err := m.conn.Query(
		*ctx, "select "+defOrderData+","+defUserData+
			" from orders o"+
			" left join users u on u.id = o.client"+
			" where o.client=$1"+
			" order by o.createddate desc",
		clientID)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"GetOrdersByClient->m.conn.Query %w", err)
	}

	defer rows.Close()

	orders := make([]ordermodel.Order, 0)
	errors := make([]error, 0)

	for rows.Next() {
		order := &ordermodel.Order{}
		user := &usermodel.User{}
		err = rows.Scan(&outOrderID, &outOrderIdentifier,
			&outOrderCreateddate, &outOrderStatus, &outOrderAccrual,
			&outOrderPointsWriteOff, &outUserID,
			&outUserLogin, &outUserPass, &outUserCreateddate)

		if err != nil {
			errors = append(errors, err)
		} else {
			user.SetUser(*outUserID, outUserLogin,
				outUserPass, outUserCreateddate)
			order.SetOrder(
				*outOrderID, outOrderIdentifier, user,
				outOrderCreateddate, outOrderStatus,
				outOrderAccrual, outOrderPointsWriteOff)

			orders = append(orders, *order)
		}
	}

	return &orders, &errors, nil
}

func (m *OrderStorage) GetOrdersByStatuses(
	ctx *context.Context,
	statuses string,
) (*[]ordermodel.Order, *[]error, error) {
	var (
		outOrderID, outUserID                   *int32
		outOrderStatus, outOrderIdentifier      *string
		outUserLogin, outUserPass               *string
		outUserCreateddate, outOrderCreateddate *time.Time
		outOrderPointsWriteOff, outOrderAccrual *float32
	)

	rows, err := m.conn.Query(
		*ctx, "select "+defOrderData+","+defUserData+
			" from orders o"+
			" left join users u on u.id = o.client"+
			" where o.status in ("+statuses+")")
	if err != nil {
		return nil, nil, fmt.Errorf(
			"GetOrdersByStatuses->m.conn.Query %w", err)
	}

	defer rows.Close()

	orders := make([]ordermodel.Order, 0)
	errors := make([]error, 0)

	for rows.Next() {
		order := &ordermodel.Order{}
		user := &usermodel.User{}
		err = rows.Scan(&outOrderID, &outOrderIdentifier,
			&outOrderCreateddate, &outOrderStatus, &outOrderAccrual,
			&outOrderPointsWriteOff, &outUserID,
			&outUserLogin, &outUserPass, &outUserCreateddate)

		if err != nil {
			errors = append(errors, err)
		} else {
			user.SetUser(*outUserID, outUserLogin,
				outUserPass, outUserCreateddate)
			order.SetOrder(
				*outOrderID, outOrderIdentifier, user,
				outOrderCreateddate, outOrderStatus,
				outOrderAccrual, outOrderPointsWriteOff)

			orders = append(orders, *order)
		}
	}

	return &orders, &errors, nil
}
