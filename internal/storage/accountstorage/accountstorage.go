package accountstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accounthistorymodel"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accountmodel"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/ordermodel"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
	"github.com/jackc/pgx/v5"
)

type AccountStorage struct {
	conn *pgx.Conn
}

func (m *AccountStorage) Initiate(
	conn *pgx.Conn,
) {
	m.conn = conn
}

const defUserData = "u.id,u.login,u.password,u.createddate"

const defAccountData = "a.id,a.points,a.withdrawn," +
	"a.client,a.createddate"

const defOrderData = "o.id, o.identifier, o.createddate, " +
	"o.status, o.accrual, o.points_write_off"

const defAccHistData = "ah.id,ah.points_write_off," +
	"ah.order,ah.createddate"

func (m *AccountStorage) MinusPointsByID(
	ctx *context.Context,
	accID int32,
	newValuePoints float32,
) (bool, error) {
	rows, err := m.conn.Exec(
		*ctx,
		"UPDATE account SET points=points-$1 where id=$2",
		newValuePoints,
		accID)
	if err != nil {
		return false, fmt.Errorf(
			"MinusPointsByID->m.conn.Exec( %w", err)
	}

	if rows.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}

func (m *AccountStorage) PlusWithdrawnByID(
	ctx *context.Context,
	accID int32,
	newValueWithdrawn float32,
) (bool, error) {
	rows, err := m.conn.Exec(
		*ctx,
		"UPDATE account SET withdrawn=withdrawn+$1 where id=$2",
		newValueWithdrawn,
		accID)
	if err != nil {
		return false, fmt.Errorf(
			"PlusWithdrawnByID->m.conn.Exec( %w", err)
	}

	if rows.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}

func (m *AccountStorage) CreateAccount(
	ctx *context.Context,
	account *accountmodel.Account,
) error {
	_, err := m.conn.Exec(
		*ctx,
		"INSERT INTO account (client,"+
			" points,withdrawn) VALUES ($1,$2,$3)",
		account.GetClient().GetID(), account.GetPoints(),
		account.GetWithdrawn())
	if err != nil {
		return fmt.Errorf(
			"CreateAccount->INSERT INTO error: %w", err)
	}

	return nil
}

func (m *AccountStorage) CreateAccountHistory(
	ctx *context.Context,
	accHist *accounthistorymodel.AccountHistory,
) error {
	_, err := m.conn.Exec(
		*ctx,
		"INSERT INTO account_history (points_write_off,order)"+
			" VALUES ($1,$2)",
		accHist.GetpointsWriteOff(), accHist.GetOrder().GetID())
	if err != nil {
		return fmt.Errorf(
			"CreateAccountHistory->INSERT INTO error: %w", err)
	}

	return nil
}

func (m *AccountStorage) GetAccountByClient(
	ctx *context.Context,
	clientID int32,
) (*accountmodel.Account, error) {
	var (
		outAccountID, outAccountClientID, outUserID int32
		outAccountCreateddate, outUserCreateddate   time.Time
		outUserLogin, outUserPass                   string
		outAccountPoints, outAccountWithdrawn       float32
	)

	err := m.conn.QueryRow(
		*ctx, "select "+defAccountData+","+defUserData+
			" from account a"+
			" left join user u on u.id = a.client"+
			" where a.client=$1"+
			" LIMIT 1",
		clientID).Scan(&outAccountID, &outAccountPoints,
		&outAccountWithdrawn, &outAccountClientID,
		&outAccountCreateddate, &outUserID,
		&outUserLogin, &outUserPass, &outUserCreateddate)
	if err != nil {
		return nil, fmt.Errorf(
			"GetAccountByClient->m.conn.QueryRow %w", err)
	}

	acc := &accountmodel.Account{}
	user := &usermodel.User{}

	user.SetUser(outUserID, outUserPass,
		outUserLogin, outUserCreateddate)
	acc.SetAccount(
		outAccountID, user, outAccountCreateddate,
		outAccountPoints, outAccountWithdrawn)

	return acc, nil
}

func (m *AccountStorage) GetAccountHistoryByClient(
	ctx *context.Context,
	clientID int32,
) (*[]accounthistorymodel.AccountHistory, *[]error, error) {
	var (
		OrdCreateddate, AccHistCreateddate          time.Time
		OrdID, AcctHistID, AccHistOrdID, OrdAccrual int32
		OrdStatus, OrdIdentifier                    string
		OrdPointsWriteOff, AccHistPointsWriteOff    float32
	)

	rows, err := m.conn.Query(
		*ctx, "select "+defOrderData+","+defAccHistData+
			" from account_history ah"+
			" left join order o on o.id = ah.order"+
			" where o.client=$1"+
			" order by ah.createddate desc",
		clientID)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"GetAccountHistoryByClient->m.conn.Query %w", err)
	}

	defer rows.Close()

	cnt := 0
	for rows.Next() {
		cnt++
	}

	accHists := make(
		[]accounthistorymodel.AccountHistory, 0, cnt)
	errors := make([]error, 0, cnt)

	if cnt == 0 {
		return &accHists, &errors, nil
	}

	for rows.Next() {
		accHist := &accounthistorymodel.AccountHistory{}
		order := &ordermodel.Order{}

		err = rows.Scan(&OrdID, &OrdIdentifier,
			&OrdCreateddate, &OrdAccrual,
			&OrdPointsWriteOff, &OrdStatus,
			&AcctHistID, &AccHistPointsWriteOff,
			&AccHistOrdID, &AccHistCreateddate)
		if err != nil {
			errors = append(errors, err)
		} else {
			order.SetOrder(OrdID, OrdIdentifier, nil,
				OrdCreateddate, OrdStatus,
				OrdAccrual, OrdPointsWriteOff)
			accHist.SetAccountHistory(AcctHistID, order,
				AccHistCreateddate, AccHistPointsWriteOff)

			accHists = append(accHists, *accHist)
		}
	}

	return &accHists, &errors, nil
}
