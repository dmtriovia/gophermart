package ordermodel

import (
	"time"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
)

const OrderStatusNew string = "NEW"

const OrderStatusProcessing string = "PROCESSING"

const OrderStatusInvalid string = "INVALID"

const OrderStatusProcessed string = "PROCESSED"

type Order struct {
	id          int32
	identifier  string
	client      *usermodel.User
	createddate time.Time
	status      string
}

func (o *Order) SetOrder(
	idDB int32,
	identifier string,
	client *usermodel.User,
	createddate time.Time,
	status string,
) {
	o.id = idDB
	o.identifier = identifier
	o.client = client
	o.createddate = createddate
	o.status = status
}

func (o *Order) SetStatus(status string) {
	o.status = status
}

func (o *Order) SetClient(user *usermodel.User) {
	o.client = user
}

func (o *Order) SetIdentifier(ident string) {
	o.identifier = ident
}

func (o *Order) GetID() int32 {
	return o.id
}

func (o *Order) GetIdentifier() string {
	return o.identifier
}

func (o *Order) GetClient() *usermodel.User {
	return o.client
}

func (o *Order) GetCreateddate() time.Time {
	return o.createddate
}
