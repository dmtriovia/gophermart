package middlewareattr

import (
	as "github.com/dmitrovia/gophermart/internal/service/authservice"
	"go.uber.org/zap"
)

type AuthMiddlewareAttr struct {
	secret      string
	zapLogger   *zap.Logger
	authService *as.AuthService
}

func (p *AuthMiddlewareAttr) Init(logger *zap.Logger,
	authService *as.AuthService,
) {
	p.secret = "qwerty"
	p.zapLogger = logger
	p.authService = authService
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
