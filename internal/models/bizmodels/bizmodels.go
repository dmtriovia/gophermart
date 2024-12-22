package bizmodels

type User struct {
	Login    string `json:"login"`
	Password string `json:"-"`
}
