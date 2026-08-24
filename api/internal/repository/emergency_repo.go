package repository

import (
	"context"

	"employee-directory-api/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmergencyContactRepository struct{ pool *pgxpool.Pool }

func NewEmergencyContactRepository(pool *pgxpool.Pool) *EmergencyContactRepository {
	return &EmergencyContactRepository{pool: pool}
}

func (r *EmergencyContactRepository) GetByPerson(ctx context.Context, personID uuid.UUID) (*domain.EmergencyContact, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, person_id, name, phone, email, relationship, created_at, updated_at FROM emergency_contacts WHERE person_id=$1`, personID)
	var ec domain.EmergencyContact
	err := row.Scan(&ec.ID, &ec.PersonID, &ec.Name, &ec.Phone, &ec.Email, &ec.Relationship, &ec.CreatedAt, &ec.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ec, nil
}

func (r *EmergencyContactRepository) Upsert(ctx context.Context, ec *domain.EmergencyContact) (*domain.EmergencyContact, error) {
	_, err := r.pool.Exec(ctx, `INSERT INTO emergency_contacts (id, person_id, name, phone, email, relationship) VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (person_id) DO UPDATE SET name=EXCLUDED.name, phone=EXCLUDED.phone, email=EXCLUDED.email, relationship=EXCLUDED.relationship, updated_at=now()`,
		ec.ID, ec.PersonID, ec.Name, ec.Phone, ec.Email, ec.Relationship)
	if err != nil {
		return nil, err
	}
	return r.GetByPerson(ctx, ec.PersonID)
}

func (r *EmergencyContactRepository) Delete(ctx context.Context, personID uuid.UUID) (bool, error) {
	result, err := r.pool.Exec(ctx, `DELETE FROM emergency_contacts WHERE person_id=$1`, personID)
	return result.RowsAffected() > 0, err
}
