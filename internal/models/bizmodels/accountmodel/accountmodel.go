package accountmodel

import (
	"time"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
)

type Account struct {
	id          int32
	createddate time.Time
	client      *usermodel.User
	points      float32
	withdrawn   float32
}

func (u *Account) SetAccount(
	idDB int32,
	client *usermodel.User,
	createddate time.Time,
	points float32,
	withdrawn float32,
) {
	u.id = idDB
	u.createddate = createddate
	u.client = client
	u.points = points
	u.withdrawn = withdrawn
}

func (u *Account) SetClient(client *usermodel.User) {
	u.client = client
}

func (u *Account) GetID() int32 {
	return u.id
}

func (u *Account) SetID(id int32) {
	u.id = id
}

func (u *Account) GetCreateddate() time.Time {
	return u.createddate
}

func (u *Account) GetWithdrawn() float32 {
	return u.withdrawn
}

func (u *Account) GetPoints() float32 {
	return u.points
}

func (u *Account) SetPoints(points float32) {
	u.points = points
}

func (u *Account) SetWithdrawn(withdrawn float32) {
	u.withdrawn = withdrawn
}

func (u *Account) GetClient() *usermodel.User {
	return u.client
}
