package accounthistorymodel

import (
	"time"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels/ordermodel"
)

type AccountHistory struct {
	id             int32
	order          *ordermodel.Order
	pointsWriteOff float32
	createddate    time.Time
}

func (u *AccountHistory) SetAccountHistory(
	idDB int32,
	order *ordermodel.Order,
	createddate time.Time,
	pointsWriteOff float32,
) {
	u.id = idDB
	u.createddate = createddate
	u.order = order
	u.pointsWriteOff = pointsWriteOff
}

func (u *AccountHistory) GetID() int32 {
	return u.id
}

func (u *AccountHistory) GetCreateddate() time.Time {
	return u.createddate
}

func (u *AccountHistory) GetpointsWriteOff() float32 {
	return u.pointsWriteOff
}

func (u *AccountHistory) GetOrder() *ordermodel.Order {
	return u.order
}

func (u *AccountHistory) SetpointsWriteOff(points float32) {
	u.pointsWriteOff = points
}

func (u *AccountHistory) SetOrder(order *ordermodel.Order) {
	u.order = order
}
