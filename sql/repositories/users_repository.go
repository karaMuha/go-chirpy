package repositories

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/karaMuha/go-chirpy/models"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) UsersRepository {
	return UsersRepository{
		db: db,
	}
}

func (r *UsersRepository) CreateUser(ctx context.Context, email string) (*models.User, *models.ResponseErr) {
	query := `
		INSERT INTO users (id, created_at, updated_at, email)
		VALUES (gen_random_uuid(), now(), now(), $1)
		RETURNING *;
	`
	row := r.db.QueryRowContext(ctx, query, email)
	var user models.User
	if err := row.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Email,
	); err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, &models.ResponseErr{
				Error:      "Email already exists",
				StatusCode: http.StatusConflict,
			}
		}
		return nil, &models.ResponseErr{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &user, nil
}

func (r *UsersRepository) ResetTable(ctx context.Context) *models.ResponseErr {
	query := `
		DELETE FROM users
	`
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return &models.ResponseErr{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}
