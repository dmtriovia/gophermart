package withdrawalsattr

import (
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
	"go.uber.org/zap"
)

type WithdrawalsAttr struct {
	zapLogger   *zap.Logger
	sessionUser *usermodel.User
}

func (p *WithdrawalsAttr) Init(
	logger *zap.Logger,
	user *usermodel.User,
) {
	p.zapLogger = logger
	p.sessionUser = user
}

func (p *WithdrawalsAttr) GetLogger() *zap.Logger {
	return p.zapLogger
}

func (p *WithdrawalsAttr) GetSessionUser() *usermodel.User {
	return p.sessionUser
}
