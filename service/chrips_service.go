package service

import (
	"context"

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
