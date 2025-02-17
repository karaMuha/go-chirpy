package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn * time.Second)),
		Subject:   userID.String(),
	})

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.UUID{}, err
	}

	if !parsedToken.Valid {
		return uuid.UUID{}, errors.New("invalid token")
	}

	userID, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	return uuid.MustParse(userID), nil
}

func GetBearerToken(headers http.Header) (string, error) {
	value := headers.Get("Authorization")
	token, found := strings.CutPrefix(value, "Bearer ")
	if !found {
		return "", errors.New("token not found")
	}

	return token, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	header := headers.Get("Authorization")
	key, found := strings.CutPrefix(header, "ApiKey ")
	if !found {
		return "", errors.New("api key not present")
	}
	return key, nil
}

func MakeRefreshToken() (string, error) {
	data := make([]byte, 32)
	_, err := rand.Read(data)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(data), nil
}
