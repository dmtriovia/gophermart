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
	var lastInsertID int32

	err := m.conn.QueryRow(
		*ctx,
		"INSERT INTO orders (identifier,client,"+
			" accrual,status) VALUES ($1,$2,$3,$4) RETURNING id",
		order.GetIdentifier(), order.GetClient().GetID(),
		order.GetAccrual(), order.GetStatus()).Scan(&lastInsertID)
	if err != nil {
		return fmt.Errorf(
			"CreateOrder->Scan: %w", err)
	}

	order.SetID(lastInsertID)

	return nil
}

func (m *OrderStorage) PlusPointsWriteOffByID(
	ctx *context.Context,
	orderID int32,
	newValuePointsWriteOff float32,
) (bool, error) {
	t := "points_write_off"

	rows, err := m.conn.Exec(
		*ctx,
		"UPDATE orders SET "+t+"="+t+"+$1 where id=$2",
		newValuePointsWriteOff,
		orderID)
	if err != nil {
		return false, fmt.Errorf(
			"PlusPointsWriteOffByID->m.conn.Exec( %w", err)
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
		outOrderID, outOrderAccrual, outUserID  int32
		outOrderStatus, outOrderIdentifier      string
		outOrderCreateddate, outUserCreateddate time.Time
		outOrderPointsWriteOff                  float32
		outUserLogin, outUserPass               string
	)

	err := m.conn.QueryRow(
		*ctx, "select "+defOrderData+","+defUserData+
			" from orders o"+
			" left join users u on u.id = order.client"+
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
		outOrderStatus, outOrderAccrual,
		outOrderPointsWriteOff)

	return order, nil
}

func (m *OrderStorage) GetOrdersByClient(
	ctx *context.Context,
	clientID int32,
) (*[]ordermodel.Order, *[]error, error) {
	var (
		outOrderID, outUserID, outOrderAccrual  int32
		outOrderStatus, outOrderIdentifier      string
		outUserLogin, outUserPass               string
		outUserCreateddate, outOrderCreateddate time.Time
		outOrderPointsWriteOff                  float32
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
			&outOrderCreateddate, &outOrderAccrual,
			&outOrderPointsWriteOff, &outOrderStatus, &outUserID,
			&outUserLogin, &outUserPass, &outUserCreateddate)

		if err != nil {
			errors = append(errors, err)
		} else {
			user.SetUser(outUserID, outUserPass,
				outUserLogin, outUserCreateddate)
			order.SetOrder(
				outOrderID, outOrderIdentifier, user,
				outOrderCreateddate, outOrderStatus,
				outOrderAccrual, outOrderPointsWriteOff)

			orders = append(orders, *order)
		}
	}

	return &orders, &errors, nil
}
