package usermodel

import "time"

type User struct {
	id          int32
	login       *string
	password    *string
	createddate *time.Time
}

func (u *User) SetUser(
	idDB int32,
	login *string,
	password *string,
	createddate *time.Time,
) {
	u.id = idDB
	u.login = login
	u.password = password
	u.createddate = createddate
}

func (u *User) SetLogin(
	login *string,
) {
	u.login = login
}

func (u *User) SetPassword(
	password *string,
) {
	u.password = password
}

func (u *User) GetID() int32 {
	return u.id
}

func (u *User) SetID(id int32) {
	u.id = id
}

func (u *User) GetLogin() *string {
	return u.login
}

func (u *User) GetPassword() *string {
	return u.password
}

func (u *User) GetCreateddate() *time.Time {
	return u.createddate
}
