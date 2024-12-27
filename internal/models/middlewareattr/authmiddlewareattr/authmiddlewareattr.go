package authmiddlewareattr

import (
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
	as "github.com/dmitrovia/gophermart/internal/service/authservice"
	"go.uber.org/zap"
)

type AuthMiddlewareAttr struct {
	secret      string
	zapLogger   *zap.Logger
	authService *as.AuthService
	sessionUser *usermodel.User
}

func (p *AuthMiddlewareAttr) Init(logger *zap.Logger,
	authService *as.AuthService,
	user *usermodel.User,
) {
	p.secret = "qwerty"
	p.zapLogger = logger
	p.authService = authService
	p.sessionUser = user
}

func (
	p *AuthMiddlewareAttr) GetSessionUser() *usermodel.User {
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
