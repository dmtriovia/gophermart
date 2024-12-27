package service

import (
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/ordermodel"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
)

type AccountService interface{}

type AuthService interface {
	UserIsExist(login string) (bool, *usermodel.User, error)
	CreateUser(user *usermodel.User) error
}

type OrderService interface {
	OrderIsExist(
		identifier string) (bool, *ordermodel.Order, error)
	CreateOrder(order *ordermodel.Order) error
}
