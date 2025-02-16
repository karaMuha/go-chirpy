package repositories

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/karaMuha/go-chirpy/models"
)

type ChirpsRepository struct {
	db *sql.DB
}

func NewChirpsRepository(db *sql.DB) ChirpsRepository {
	return ChirpsRepository{
		db: db,
	}
}

func (r *ChirpsRepository) CreateChirp(ctx context.Context, body, userID string) (*models.Chirp, *models.ResponseErr) {
	query := `
		INSERT INTO chirps (id, created_at, updated_at, body, user_id)
		VALUES (gen_random_uuid (), now(), now(), $1, $2)
		RETURNING *;
	`
	row := r.db.QueryRowContext(ctx, query, body, userID)

	var chirp models.Chirp
	if err := row.Scan(
		&chirp.ID,
		&chirp.CreatedAt,
		&chirp.UpdatedAt,
		&chirp.Body,
		&chirp.UserID,
	); err != nil {
		return nil, &models.ResponseErr{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &chirp, nil
}
