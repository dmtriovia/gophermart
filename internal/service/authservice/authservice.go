package authservice

import "github.com/dmitrovia/gophermart/internal/storage"

type AuthService struct {
	repository storage.Storage
}

func NewAuthService(
	stor storage.Storage,
) *AuthService {
	return &AuthService{repository: stor}
}
