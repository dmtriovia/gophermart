package storage

import (
	"context"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels"
)

type Storage interface {
	GetUser(
		ctx *context.Context,
		login string) (*bizmodels.User, error)

	CreateUser(
		ctx *context.Context,
		user *bizmodels.User) error
}
