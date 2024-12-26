package bizmodels

import "time"

type User struct {
	id          int32
	login       string
	password    string
	createddate time.Time
}

func (u *User) SetUser(
	id int32,
	login string,
	password string,
	createddate time.Time,
) {
	u.id = id
	u.login = login
	u.password = password
	u.createddate = createddate
}

func (u *User) SetLogin(
	login string,
) {
	u.login = login
}

func (u *User) SetPassword(
	password string,
) {
	u.password = password
}

func (u *User) GetID() int32 {
	return u.id
}

func (u *User) GetLogin() string {
	return u.login
}

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) GetCreateddate() time.Time {
	return u.createddate
}

type Order struct {
	id          int32
	identifier  string
	client      *User
	createddate time.Time
}

func (o *Order) SetOrder(
	id int32,
	identifier string,
	client *User,
	createddate time.Time,
) {
	o.id = id
	o.identifier = identifier
	o.client = client
	o.createddate = createddate
}

func (o *Order) SetClient(user *User) {
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

func (o *Order) GetClient() *User {
	return o.client
}

func (o *Order) GetCreateddate() time.Time {
	return o.createddate
}
