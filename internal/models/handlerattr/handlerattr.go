package handlerattr

type LoginAttr struct {
	Secret       string
	TokenExpHour int
}

type RegisterAttr struct {
	Secret       string
	TokenExpHour int
}
