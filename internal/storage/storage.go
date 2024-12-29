package storage

import (
	"context"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accounthistorymodel"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accountmodel"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/ordermodel"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
)

type UserStorage interface {
	GetUser(
		ctx *context.Context,
		login string) (*usermodel.User, error)

	CreateUser(
		ctx *context.Context,
		user *usermodel.User) error
}

type OrderStorage interface {
	GetOrder(
		ctx *context.Context,
		ident string) (*ordermodel.Order, error)

	GetOrdersByClient(
		ctx *context.Context,
		clientID int32) (*[]ordermodel.Order, *[]error, error)

	CreateOrder(
		ctx *context.Context,
		user *ordermodel.Order) error

	PlusPointsWriteOffByID(
		ctx *context.Context,
		orderID int32,
		newValuePointsWriteOff float32,
	) (bool, error)
}

type AccountStorage interface {
	CreateAccount(ctx *context.Context,
		account *accountmodel.Account) error

	GetAccountByClient(
		ctx *context.Context,
		clientID int32,
	) (*accountmodel.Account, error)

	PlusWithdrawnByID(
		ctx *context.Context,
		accID int32,
		newValueWithdrawn float32,
	) (bool, error)

	MinusPointsByID(
		ctx *context.Context,
		accID int32,
		newValuePoints float32,
	) (bool, error)

	CreateAccountHistory(ctx *context.Context,
		account *accounthistorymodel.AccountHistory) error
}
