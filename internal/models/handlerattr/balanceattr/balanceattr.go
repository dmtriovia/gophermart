package balanceattr

import (
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
	"go.uber.org/zap"
)

type BalanceAttr struct {
	zapLogger   *zap.Logger
	sessionUser *usermodel.User
}

func (p *BalanceAttr) Init(
	logger *zap.Logger,
	user *usermodel.User,
) {
	p.zapLogger = logger
	p.sessionUser = user
}

func (p *BalanceAttr) GetLogger() *zap.Logger {
	return p.zapLogger
}

func (p *BalanceAttr) GetSessionUser() *usermodel.User {
	return p.sessionUser
}
