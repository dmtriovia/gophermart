package loginattr

import "go.uber.org/zap"

type LoginAttr struct {
	secret       string
	tokenExpHour int
	zapLogger    *zap.Logger
}

func (p *LoginAttr) Init(logger *zap.Logger) {
	p.secret = "qwerty"
	p.tokenExpHour = 24
	p.zapLogger = logger
}

func (p *LoginAttr) GetSecret() string {
	return p.secret
}

func (p *LoginAttr) GetTokenExpHour() int {
	return p.tokenExpHour
}

func (p *LoginAttr) GetLogger() *zap.Logger {
	return p.zapLogger
}
