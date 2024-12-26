package authservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels"
	"github.com/dmitrovia/gophermart/internal/storage"
)

type AuthService struct {
	repository  storage.Storage
	ctxDuration time.Duration
}

func NewAuthService(
	stor storage.Storage, ctxDur int,
) *AuthService {
	return &AuthService{
		repository:  stor,
		ctxDuration: time.Duration(ctxDur),
	}
}

func (s *AuthService) UserIsExist(login string) (
	bool, *bizmodels.User, error,
) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		s.ctxDuration)
	defer cancel()

	user, err := s.repository.GetUser(&ctx, login)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil, nil
	}

	if err != nil {
		return false, nil, fmt.Errorf(
			"UserIsExist->GetUser: %w",
			err)
	}

	return true, user, nil
}

func (s *AuthService) CreateUser(
	user *bizmodels.User,
) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		s.ctxDuration)
	defer cancel()

	err := s.repository.CreateUser(&ctx, user)
	if err != nil {
		return fmt.Errorf(
			"CreateUser->s.repository.CreateUser: %w",
			err)
	}

	return nil
}
