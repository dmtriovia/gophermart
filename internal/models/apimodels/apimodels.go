package apimodels

type RegisterUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SetOrder struct {
	Identifier string `json:"identifier"`
}
