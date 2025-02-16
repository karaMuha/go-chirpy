package service

import (
	"context"
	"net/http"

	"github.com/karaMuha/go-chirpy/models"
	"github.com/karaMuha/go-chirpy/sql/repositories"
)

type ChirpsService struct {
	chripRepo repositories.ChirpsRepository
}

func NewChripsService(chirpRepo repositories.ChirpsRepository) ChirpsService {
	return ChirpsService{
		chripRepo: chirpRepo,
	}
}

func (s *ChirpsService) CreateChrip(ctx context.Context, body, userID string) (*models.Chirp, *models.ResponseErr) {
	return s.chripRepo.CreateChirp(ctx, body, userID)
}

func (s *ChirpsService) GetAll(ctx context.Context, authorID, sorting string) (*[]models.Chirp, *models.ResponseErr) {
	return s.chripRepo.GetAll(ctx, authorID, sorting)
}

func (s *ChirpsService) GetByID(ctx context.Context, chirpID string) (*models.Chirp, *models.ResponseErr) {
	return s.chripRepo.GetChirpByID(ctx, chirpID)
}

func (s *ChirpsService) Delete(ctx context.Context, userID, chirpID string) *models.ResponseErr {
	chirp, respErr := s.chripRepo.GetChirpByID(ctx, chirpID)
	if respErr != nil {
		return respErr
	}

	if chirp.UserID.String() != userID {
		return &models.ResponseErr{
			Error:      "Not your chirp",
			StatusCode: http.StatusForbidden,
		}
	}

	return s.chripRepo.DeleteChirp(ctx, chirpID)
}
