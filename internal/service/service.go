package service

import (
	"context"

	"github.com/dunkykorZhik/avito-tech/internal/entity"
	"github.com/dunkykorZhik/avito-tech/internal/repo"
)

type User interface {
	GetUser(ctx context.Context, id int64) (*entity.User, error)
}

type Service struct {
	User User
}

type ServiceDependencies struct {
	Repo *repo.Repositories
}

func NewService(deps ServiceDependencies) *Service {
	return &Service{
		User: NewUserService(deps.Repo.User),
	}
}
