package getordersattr

import (
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
	"go.uber.org/zap"
)

type GetOrdersAttr struct {
	zapLogger   *zap.Logger
	sessionUser *usermodel.User
}

func (p *GetOrdersAttr) Init(
	logger *zap.Logger,
	user *usermodel.User,
) {
	p.zapLogger = logger
	p.sessionUser = user
}

func (p *GetOrdersAttr) GetLogger() *zap.Logger {
	return p.zapLogger
}

func (p *GetOrdersAttr) GetSessionUser() *usermodel.User {
	return p.sessionUser
}
