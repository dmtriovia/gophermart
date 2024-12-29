package service

import (
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accountmodel"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/ordermodel"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
)

type AccountService interface {
	CreateAccount(account *accountmodel.Account) error
	GetAccountByClient(clientID int32,
	) (*accountmodel.Account, error)

	CalculatePoints(
		acc *accountmodel.Account,
		order *ordermodel.Order,
		points float32,
	) error
}

type AuthService interface {
	UserIsExist(login string) (bool, *usermodel.User, error)
	CreateUser(user *usermodel.User) error
}

type OrderService interface {
	OrderIsExist(
		identifier string) (bool, *ordermodel.Order, error)
	CreateOrder(order *ordermodel.Order) error
	GetOrdersByClient(
		clientID int32) (*[]ordermodel.Order, *[]error, error)
}
