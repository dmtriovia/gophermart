package apimodels

import "time"

type InRegisterUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type InLoginUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type InSetOrder struct {
	Identifier string `json:"identifier"`
}

type OutGetOrder struct {
	Identifier  string    `json:"number"`
	Createddate time.Time `json:"uploaded_at"`
	Status      string    `json:"status"`
	Accrual     int32     `json:"accrual,omitempty"`
}
