package repository

import (
	"context"

	"employee-directory-api/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrgRepository struct{ pool *pgxpool.Pool }

func NewOrgRepository(pool *pgxpool.Pool) *OrgRepository { return &OrgRepository{pool: pool} }

func (r *OrgRepository) Create(ctx context.Context, o *domain.Organization) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO organizations (id, name, type, country, timezone) VALUES ($1,$2,$3,$4,$5)`, o.ID, o.Name, o.Type, o.Country, o.Timezone)
	return err
}

func (r *OrgRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, name, type, country, timezone, is_active, created_at, updated_at FROM organizations WHERE id=$1`, id)
	var o domain.Organization
	err := row.Scan(&o.ID, &o.Name, &o.Type, &o.Country, &o.Timezone, &o.IsActive, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &o, nil
}

func (r *OrgRepository) List(ctx context.Context) ([]domain.Organization, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, type, country, timezone, is_active, created_at, updated_at FROM organizations ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Organization
	for rows.Next() {
		var o domain.Organization
		if err := rows.Scan(&o.ID, &o.Name, &o.Type, &o.Country, &o.Timezone, &o.IsActive, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}
