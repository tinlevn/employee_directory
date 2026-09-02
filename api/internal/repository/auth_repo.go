package repository

import (
	"context"

	"employee-directory-api/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{pool: pool}
}

func (r *AuthRepository) CreateAccount(ctx context.Context, account *domain.PersonAccount) error {
	q := `INSERT INTO person_accounts (person_id, username, password_hash, role, permissions, is_active)
	      VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	err := r.pool.QueryRow(ctx, q, account.PersonID, account.Username, account.PasswordHash, account.Role, account.Permissions, account.IsActive).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)
	return err
}

func (r *AuthRepository) GetAccountByUsername(ctx context.Context, username string) (*domain.PersonAccount, error) {
	q := `SELECT id, person_id, username, password_hash, role, permissions, is_active, last_login, two_factor_enabled, account_locked, created_at, updated_at
	      FROM person_accounts WHERE username = $1`
	var acc domain.PersonAccount
	err := r.pool.QueryRow(ctx, q, username).Scan(
		&acc.ID, &acc.PersonID, &acc.Username, &acc.PasswordHash, &acc.Role, &acc.Permissions,
		&acc.IsActive, &acc.LastLogin, &acc.TwoFactorEnabled, &acc.AccountLocked, &acc.CreatedAt, &acc.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &acc, nil
}

func (r *AuthRepository) GetOrgIDByPersonID(ctx context.Context, personID uuid.UUID) (uuid.UUID, error) {
	q := `SELECT org_id FROM persons WHERE id = $1`
	var orgID uuid.UUID
	err := r.pool.QueryRow(ctx, q, personID).Scan(&orgID)
	if err != nil {
		return uuid.Nil, err
	}
	return orgID, nil
}
