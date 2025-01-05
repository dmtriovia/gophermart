package ordermodel

import (
	"time"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
)

const OrderStatusRegistered string = "REGISTERED"

const OrderStatusProcessing string = "PROCESSING"

const OrderStatusInvalid string = "INVALID"

const OrderStatusProcessed string = "PROCESSED"

type Order struct {
	id             int32
	identifier     *string
	client         *usermodel.User
	createddate    *time.Time
	status         *string
	accrual        *float32
	pointsWriteOff *float32
}

func (o *Order) SetOrder(
	idDB int32,
	identifier *string,
	client *usermodel.User,
	createddate *time.Time,
	status *string,
	accrual *float32,
	pointsWriteOff *float32,
) {
	o.id = idDB
	o.identifier = identifier
	o.client = client
	o.createddate = createddate
	o.status = status
	o.accrual = accrual
	o.pointsWriteOff = pointsWriteOff
}

func (o *Order) GetpointsWriteOff() *float32 {
	return o.pointsWriteOff
}

func (o *Order) SetpointsWriteOff(points *float32) {
	o.pointsWriteOff = points
}

func (o *Order) SetStatus(status *string) {
	o.status = status
}

func (o *Order) SetClient(user *usermodel.User) {
	o.client = user
}

func (o *Order) SetIdentifier(ident *string) {
	o.identifier = ident
}

func (o *Order) GetID() int32 {
	return o.id
}

func (o *Order) SetID(id int32) {
	o.id = id
}

func (o *Order) GetAccrual() *float32 {
	return o.accrual
}

func (o *Order) GetStatus() *string {
	return o.status
}

func (o *Order) GetIdentifier() *string {
	return o.identifier
}

func (o *Order) GetClient() *usermodel.User {
	return o.client
}

func (o *Order) GetCreateddate() *time.Time {
	return o.createddate
}
