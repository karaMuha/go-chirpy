package repositories

import (
	"context"
	"database/sql"
	"fmt"
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

func (r *ChirpsRepository) GetAll(ctx context.Context, authorID, sorting string) (*[]models.Chirp, *models.ResponseErr) {
	var rows *sql.Rows
	var err error

	if authorID == "" {
		query := fmt.Sprintf(`
		SELECT *
		FROM chirps
		ORDER BY created_at %s
	`, sorting)
		rows, err = r.db.QueryContext(ctx, query)
	} else {
		query := fmt.Sprintf(`
		SELECT *
		FROM chirps
		WHERE user_id = $1
		ORDER BY created_at %s
	`, sorting)
		rows, err = r.db.QueryContext(ctx, authorID, query)
	}

	if err != nil {
		return nil, &models.ResponseErr{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	var chripList []models.Chirp
	for rows.Next() {
		var chirp models.Chirp
		err := rows.Scan(
			&chirp.ID,
			&chirp.CreatedAt,
			&chirp.UpdatedAt,
			&chirp.Body,
			&chirp.UserID,
		)
		if err != nil {
			return nil, &models.ResponseErr{
				Error:      err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
		}
		chripList = append(chripList, chirp)
	}

	err = rows.Err()
	if err != nil {
		return nil, &models.ResponseErr{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &chripList, nil
}

func (r *ChirpsRepository) GetChirpByID(ctx context.Context, chirpID string) (*models.Chirp, *models.ResponseErr) {
	query := `
		SELECT *
		FROM chirps
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, chirpID)
	var chirp models.Chirp
	err := row.Scan(&chirp.ID, &chirp.CreatedAt, &chirp.UpdatedAt, &chirp.Body, &chirp.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &models.ResponseErr{
				Error:      "Not found",
				StatusCode: http.StatusNotFound,
			}
		}
		return nil, &models.ResponseErr{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &chirp, nil
}

func (r *ChirpsRepository) DeleteChirp(ctx context.Context, chirpID string) *models.ResponseErr {
	query := `
		DELETE FROM chirps
		WHERE id = $1;
	`
	res, err := r.db.ExecContext(ctx, query, chirpID)
	if err != nil {
		return &models.ResponseErr{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &models.ResponseErr{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if rowsAffected == 0 {
		return &models.ResponseErr{
			Error:      "Chirp not found",
			StatusCode: http.StatusNotFound,
		}
	}

	return nil
}
