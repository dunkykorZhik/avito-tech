package service

import (
	"context"

	"github.com/dunkykorZhik/avito-tech/internal/entity"
	"github.com/dunkykorZhik/avito-tech/internal/repo"
)

type UserService struct {
	userRepo repo.User
}

func NewUserService(userRepo repo.User) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUser(ctx context.Context, username string) (*entity.User, error) {
	return s.userRepo.GetUserByName(ctx, username)
}
