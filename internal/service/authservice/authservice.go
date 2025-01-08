package authservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
	"github.com/dmitrovia/gophermart/internal/storage"
)

type AuthService struct {
	repository  storage.UserStorage
	ctxDuration time.Duration
}

func NewAuthService(
	stor storage.UserStorage,
	ctxDur time.Duration,
) *AuthService {
	return &AuthService{
		repository:  stor,
		ctxDuration: ctxDur,
	}
}

func (s *AuthService) UserIsExist(login string) (
	bool, *usermodel.User, error,
) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		s.ctxDuration)
	defer cancel()

	user, err := s.repository.GetUser(&ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil, nil
		}

		return false, nil, fmt.Errorf(
			"UserIsExist->GetUser: %w",
			err)
	}

	return true, user, nil
}

func (s *AuthService) CreateUser(
	user *usermodel.User,
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
