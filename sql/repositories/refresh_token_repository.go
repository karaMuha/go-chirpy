package repositories

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/karaMuha/go-chirpy/models"
)

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) RefreshTokenRepository {
	return RefreshTokenRepository{
		db: db,
	}
}

func (r *RefreshTokenRepository) SaveRefreshToken(ctx context.Context, token, userID string, expirationDate time.Time) *models.ResponseErr {
	query := `
		INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
		VALUES ($1, now(), now(), $2, $3);
	`
	_, err := r.db.ExecContext(ctx, query, token, userID, expirationDate)
	if err != nil {
		return &models.ResponseErr{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (r *RefreshTokenRepository) GetToken(ctx context.Context, token string) (*models.RefreshToken, *models.ResponseErr) {
	query := `
		SELECT *
		FROM refresh_tokens
		WHERE token = $1
	`
	row := r.db.QueryRowContext(ctx, query, token)

	var refreshToken models.RefreshToken
	var revokedAt sql.NullTime
	if err := row.Scan(
		&refreshToken.Token,
		&refreshToken.CreatedAt,
		&refreshToken.UpdatedAt,
		&refreshToken.UserID,
		&refreshToken.ExpiresAt,
		&revokedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, &models.ResponseErr{
				Error:      "Refresh token not found",
				StatusCode: http.StatusUnauthorized,
			}
		}
		return nil, &models.ResponseErr{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	refreshToken.RevokedAt = revokedAt

	return &refreshToken, nil
}

func (r *RefreshTokenRepository) RevokeToken(ctx context.Context, token string) *models.ResponseErr {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = now(), updated_at = now()
		WHERE token = $1;
	`
	_, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return &models.ResponseErr{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}
