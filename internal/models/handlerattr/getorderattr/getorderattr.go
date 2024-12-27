package getorderattr

import (
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
	"go.uber.org/zap"
)

type GetOrderAttr struct {
	zapLogger   *zap.Logger
	sessionUser *usermodel.User
}

func (p *GetOrderAttr) Init(
	logger *zap.Logger,
	user *usermodel.User,
) {
	p.zapLogger = logger
	p.sessionUser = user
}

func (p *GetOrderAttr) GetLogger() *zap.Logger {
	return p.zapLogger
}

func (p *GetOrderAttr) GetSessionUser() *usermodel.User {
	return p.sessionUser
}
