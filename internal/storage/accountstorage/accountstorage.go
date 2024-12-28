package accountstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accountmodel"
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
