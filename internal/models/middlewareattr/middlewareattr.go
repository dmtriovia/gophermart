package middlewareattr

import (
	"github.com/dmitrovia/gophermart/internal/models/bizmodels"
	as "github.com/dmitrovia/gophermart/internal/service/authservice"
	"go.uber.org/zap"
)

type AuthMiddlewareAttr struct {
	secret      string
	zapLogger   *zap.Logger
	authService *as.AuthService
	sessionUser *bizmodels.User
}

func (p *AuthMiddlewareAttr) Init(logger *zap.Logger,
	authService *as.AuthService,
	user *bizmodels.User,
) {
	p.secret = "qwerty"
	p.zapLogger = logger
	p.authService = authService
	p.sessionUser = user
}

func (
	p *AuthMiddlewareAttr) GetSessionUser() *bizmodels.User {
	return p.sessionUser
}

func (p *AuthMiddlewareAttr) GetSecret() string {
	return p.secret
}

func (p *AuthMiddlewareAttr) GetLogger() *zap.Logger {
	return p.zapLogger
}

func (
	p *AuthMiddlewareAttr) GetAuthService() *as.AuthService {
	return p.authService
}
