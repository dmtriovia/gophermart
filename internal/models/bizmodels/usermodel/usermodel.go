package usermodel

import "time"

type User struct {
	id          int32
	login       string
	password    string
	createddate time.Time
	points      float32
	withdrawn   float32
}

func (u *User) SetUser(
	idDB int32,
	login string,
	password string,
	createddate time.Time,
	points float32,
	withdrawn float32,
) {
	u.id = idDB
	u.login = login
	u.password = password
	u.createddate = createddate
	u.points = points
	u.withdrawn = withdrawn
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

func (u *User) GetWithdrawn() float32 {
	return u.withdrawn
}

func (u *User) GetPoints() float32 {
	return u.points
}
