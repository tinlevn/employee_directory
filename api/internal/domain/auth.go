package domain

import (
	"time"

	"github.com/google/uuid"
)

type PersonAccount struct {
	ID               uuid.UUID  `json:"id"`
	PersonID         uuid.UUID  `json:"person_id"`
	Username         string     `json:"username"`
	PasswordHash     string     `json:"-"`
	Role             string     `json:"role"`
	Permissions      []string   `json:"permissions"`
	IsActive         bool       `json:"is_active"`
	LastLogin        *time.Time `json:"last_login,omitempty"`
	TwoFactorEnabled bool       `json:"two_factor_enabled"`
	AccountLocked    bool       `json:"account_locked"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type RegisterRequest struct {
	PersonID string `json:"person_id" validate:"required,uuid"`
	Username string `json:"username" validate:"required,min=3,max=200"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"omitempty,oneof=admin manager staff read-only"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token   string         `json:"token"`
	Account *PersonAccount `json:"account"`
	OrgID   uuid.UUID      `json:"org_id"` // Included for convenience
}
