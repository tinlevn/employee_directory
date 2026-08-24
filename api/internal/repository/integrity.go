package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrOrganizationMismatch = errors.New("referenced record belongs to another organization")

func ensurePersonInOrg(ctx context.Context, pool *pgxpool.Pool, personID, orgID uuid.UUID) error {
	var actualOrg uuid.UUID
	err := pool.QueryRow(ctx, `SELECT org_id FROM persons WHERE id=$1`, personID).Scan(&actualOrg)
	if err != nil {
		return err
	}
	if actualOrg != orgID {
		return fmt.Errorf("%w: person %s", ErrOrganizationMismatch, personID)
	}
	return nil
}
