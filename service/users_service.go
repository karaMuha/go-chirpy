package service

import (
	"context"
	"net/http"
	"time"

	"github.com/karaMuha/go-chirpy/internal/auth"
	"github.com/karaMuha/go-chirpy/models"
	"github.com/karaMuha/go-chirpy/sql/repositories"
	"github.com/karaMuha/go-chirpy/state"
)

type UsersService struct {
	usersRepository repositories.UsersRepository
	appState        *state.AppState
}

func NewUsersService(usersRepository repositories.UsersRepository, appState *state.AppState) UsersService {
	return UsersService{
		usersRepository: usersRepository,
		appState:        appState,
	}
}

func (s *UsersService) CreateUser(ctx context.Context, email, password string) (*models.User, *models.ResponseErr) {
	return s.usersRepository.CreateUser(ctx, email, password)
}

func (s *UsersService) ResetUsers(ctx context.Context) *models.ResponseErr {
	return s.usersRepository.ResetTable(ctx)
}

func (s *UsersService) Login(ctx context.Context, email, password string, expirationDuration int) (*models.User, *models.ResponseErr) {
	user, respErr := s.usersRepository.GetByEmail(ctx, email)
	if respErr != nil {
		return nil, respErr
	}

	if err := auth.CheckPassword(password, user.Password); err != nil {
		return nil, &models.ResponseErr{
			Error:      "incorrect email or password",
			StatusCode: http.StatusUnauthorized,
		}
	}

	token, err := auth.MakeJWT(user.ID, s.appState.Secret, time.Duration(expirationDuration))
	if err != nil {
		return nil, &models.ResponseErr{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	user.Token = token

	return user, nil
}

func (s *UsersService) UpdateAccount(ctx context.Context, userID, email, password string) (*models.User, *models.ResponseErr) {
	user, respErr := s.usersRepository.GetByID(ctx, userID)
	if respErr != nil {
		return nil, respErr
	}
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, &models.ResponseErr{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	user.Email = email
	user.Password = hashedPassword

	updatedUser, respErr := s.usersRepository.UpdateAccount(ctx, userID, email, hashedPassword)
	return updatedUser, respErr
}

func (s *UsersService) UpgradeToRed(ctx context.Context, userID string) *models.ResponseErr {
	return s.usersRepository.UpgradeToRed(ctx, userID)
}
