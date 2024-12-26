package handlerattr

import (
	"github.com/dmitrovia/gophermart/internal/models/bizmodels"
	"go.uber.org/zap"
)

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

type RegisterAttr struct {
	secret       string
	tokenExpHour int
	zapLogger    *zap.Logger
}

func (p *RegisterAttr) Init(logger *zap.Logger) {
	p.secret = "qwerty"
	p.tokenExpHour = 24
	p.zapLogger = logger
}

func (p *RegisterAttr) GetSecret() string {
	return p.secret
}

func (p *RegisterAttr) GetTokenExpHour() int {
	return p.tokenExpHour
}

func (p *RegisterAttr) GetLogger() *zap.Logger {
	return p.zapLogger
}

type SetOrderAttr struct {
	zapLogger              *zap.Logger
	validIdentOrderPattern string
	sessionUser            *bizmodels.User
}

func (p *SetOrderAttr) Init(
	logger *zap.Logger,
	user *bizmodels.User,
) {
	p.zapLogger = logger
	p.validIdentOrderPattern = "^[0-9]+$"
	p.sessionUser = user
}

func (p *SetOrderAttr) GetLogger() *zap.Logger {
	return p.zapLogger
}

func (p *SetOrderAttr) GetSessionUser() *bizmodels.User {
	return p.sessionUser
}

func (p *SetOrderAttr) GetValidIdentOrderPattern() string {
	return p.validIdentOrderPattern
}
