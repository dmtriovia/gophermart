package accountstorage

import (
	"context"
	"fmt"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accountmodel"
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
