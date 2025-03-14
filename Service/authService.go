package service

import (
	"context"
	"errors"
	"redditBack/model"
	"redditBack/repository"
	"time"
)

type AuthService struct {
	userRepo  repository.UserRepository
	cacheRepo repository.CacheRepository
}

func NewAuthService(userRepo repository.UserRepository, cacheRepo repository.CacheRepository) AuthService {
	return AuthService{userRepo: userRepo,
		cacheRepo: cacheRepo}
}

func (s *AuthService) Register(ctx context.Context, user *model.User) error {

	existingUser, err := s.userRepo.FindByUsername(ctx, user.Username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("username already exists")
	}
	existingUser2, _ := s.userRepo.FindByEmail(ctx, user.Email)

	if existingUser2 != nil {
		return errors.New("username already exists")
	}

	return s.userRepo.Create(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*model.User, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	if user.PasswordHash != password {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *AuthService) InvalidateToken(ctx context.Context, tokenString string) error {

	expiration := 24 * time.Hour
	return s.cacheRepo.InvalidateToken(ctx, tokenString, expiration)
}

func (s *AuthService) IsTokenValid(ctx context.Context, tokenString string) (bool, error) {
	return s.cacheRepo.IsTokenInvalid(ctx, tokenString)
}
