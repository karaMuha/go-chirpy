package service

import (
	"net/http"
	"strings"

	"github.com/karaMuha/go-chirpy/models"
)

type Service struct {
	profane map[string]string
}

func NewService() Service {
	profane := make(map[string]string)
	profane["kerfuffle"] = "kerfuffle"
	profane["sharbert"] = "sharbert"
	profane["fornax"] = "fornax"

	return Service{
		profane: profane,
	}
}

type Response struct {
	CleanedBody string `json:"cleaned_body,omitempty"`
}

func (s *Service) ValidateChirp(chirp models.Chirp) (*Response, *models.ResponseErr) {
	if len(chirp.Body) > 140 {
		respErr := models.ResponseErr{
			Error:      "Chirp is too long",
			StatusCode: http.StatusBadRequest,
		}
		return nil, &respErr
	}

	words := strings.Fields(chirp.Body)
	for i, word := range words {
		if _, ok := s.profane[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}

	response := Response{CleanedBody: strings.Join(words, " ")}
	return &response, nil
}
