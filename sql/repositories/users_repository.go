package repositories

import (
	"context"
	"database/sql"
	"fmt"
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

func (r *UsersRepository) CreateUser(ctx context.Context, email, password string) (*models.User, *models.ResponseErr) {
	if r.db == nil {
		fmt.Println("UserRepo DB is nil")
	}
	query := `
		INSERT INTO users (id, created_at, updated_at, email, hashed_password)
		VALUES (gen_random_uuid(), now(), now(), $1, $2)
		RETURNING *;
	`
	row := r.db.QueryRowContext(ctx, query, email, password)
	var user models.User
	if err := row.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Email,
		&user.Password,
		&user.IsChirpyRed,
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

func (r *UsersRepository) GetByID(ctx context.Context, userID string) (*models.User, *models.ResponseErr) {
	query := `
		SELECT *
		FROM users
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, userID)
	var user models.User
	if err := row.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Email,
		&user.Password,
		&user.IsChirpyRed,
	); err != nil {
		return nil, &models.ResponseErr{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &user, nil
}

func (r *UsersRepository) GetByEmail(ctx context.Context, email string) (*models.User, *models.ResponseErr) {
	query := `
		SELECT *
		FROM users
		WHERE email = $1
	`
	row := r.db.QueryRowContext(ctx, query, email)

	var user models.User
	if err := row.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Email,
		&user.Password,
		&user.IsChirpyRed,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, &models.ResponseErr{
				Error:      "User not found",
				StatusCode: http.StatusNotFound,
			}
		}
		return nil, &models.ResponseErr{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &user, nil
}

func (r *UsersRepository) UpdateAccount(ctx context.Context, userID, email, password string) (*models.User, *models.ResponseErr) {
	query := `
		UPDATE users
		SET email = $1, hashed_password = $2, updated_at = now()
		WHERE id = $3
		RETURNING *;
	`
	row := r.db.QueryRowContext(ctx, query, email, password, userID)

	var user models.User
	if err := row.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Email,
		&user.Password,
		&user.IsChirpyRed,
	); err != nil {
		return nil, &models.ResponseErr{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &user, nil
}

func (r *UsersRepository) UpgradeToRed(ctx context.Context, userID string) *models.ResponseErr {
	query := `
		UPDATE users
		SET is_chirpy_red = true
		WHERE id = $1
	`
	res, err := r.db.ExecContext(ctx, query, userID)
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
			Error:      "User not found",
			StatusCode: http.StatusNotFound,
		}
	}

	return nil
}
