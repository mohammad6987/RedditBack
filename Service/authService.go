package service

import (
	"context"
	"errors"
	"redditBack/model"
	"redditBack/repository"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, user *model.User) error {

	existingUser, err := s.userRepo.FindByUsername(ctx, user.Username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("username already exists")
	}
	existingUser2, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
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
