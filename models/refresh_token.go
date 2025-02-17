package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Token     string       `json:"refresh_token"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	UserID    uuid.UUID    `json:"user_id"`
	ExpiresAt time.Time    `json:"expires_at"`
	RevokedAt sql.NullTime `json:"revoked_at"`
}
