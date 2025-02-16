package service

import (
	"context"

	"github.com/karaMuha/go-chirpy/models"
	"github.com/karaMuha/go-chirpy/sql/repositories"
)

type UsersService struct {
	usersRepository repositories.UsersRepository
}

func NewUsersService(usersRepository repositories.UsersRepository) UsersService {
	return UsersService{
		usersRepository: usersRepository,
	}
}

func (s *UsersService) CreateUser(ctx context.Context, email string) (*models.User, *models.ResponseErr) {
	return s.usersRepository.CreateUser(ctx, email)
}

func (s *UsersService) ResetUsers(ctx context.Context) *models.ResponseErr {
	return s.usersRepository.ResetTable(ctx)
}
