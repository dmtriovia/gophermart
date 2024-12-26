package service

import "github.com/dmitrovia/gophermart/internal/models/bizmodels"

type AccountService interface{}

type AuthService interface {
	UserIsExist(login string) (bool, *bizmodels.User, error)
	CreateUser(user *bizmodels.User) error
}

type OrderService interface {
	OrderIsExist(
		identifier string) (bool, *bizmodels.Order, error)
	CreateOrder(order *bizmodels.Order) error
}
