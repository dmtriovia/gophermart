package setorderattr

import (
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
	"go.uber.org/zap"
)

type SetOrderAttr struct {
	zapLogger              *zap.Logger
	validIdentOrderPattern string
	sessionUser            *usermodel.User
}

func (p *SetOrderAttr) Init(
	logger *zap.Logger,
	user *usermodel.User,
) {
	p.zapLogger = logger
	p.validIdentOrderPattern = "^[0-9]+$"
	p.sessionUser = user
}

func (p *SetOrderAttr) GetLogger() *zap.Logger {
	return p.zapLogger
}

func (p *SetOrderAttr) GetSessionUser() *usermodel.User {
	return p.sessionUser
}

func (p *SetOrderAttr) GetValidIdentOrderPattern() string {
	return p.validIdentOrderPattern
}
