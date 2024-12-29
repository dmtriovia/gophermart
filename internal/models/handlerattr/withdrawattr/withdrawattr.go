package withdrawattr

import (
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
	"go.uber.org/zap"
)

type WithdrawAttr struct {
	zapLogger              *zap.Logger
	sessionUser            *usermodel.User
	validIdentOrderPattern string
}

func (p *WithdrawAttr) Init(
	logger *zap.Logger,
	user *usermodel.User,
) {
	p.zapLogger = logger
	p.sessionUser = user
	p.validIdentOrderPattern = "^[0-9]+$"
}

func (p *WithdrawAttr) GetLogger() *zap.Logger {
	return p.zapLogger
}

func (p *WithdrawAttr) GetSessionUser() *usermodel.User {
	return p.sessionUser
}

func (p *WithdrawAttr) GetValidIdentOrderPattern() string {
	return p.validIdentOrderPattern
}
