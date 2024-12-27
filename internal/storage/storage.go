package storage

import (
	"context"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels/ordermodel"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
)

type Storage interface {
	GetUser(
		ctx *context.Context,
		login string) (*usermodel.User, error)

	CreateUser(
		ctx *context.Context,
		user *usermodel.User) error

	GetOrder(
		ctx *context.Context,
		login string) (*ordermodel.Order, error)

	CreateOrder(
		ctx *context.Context,
		user *ordermodel.Order) error
}
