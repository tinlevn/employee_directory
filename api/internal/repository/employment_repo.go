package repository

import (
	"context"
	"fmt"

	"employee-directory-api/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmploymentRepository struct {
	pool *pgxpool.Pool
}

func NewEmploymentRepository(pool *pgxpool.Pool) *EmploymentRepository {
	return &EmploymentRepository{pool: pool}
}

func (r *EmploymentRepository) CreateVersioned(ctx context.Context, rec *domain.EmploymentRecord) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := ensurePersonInOrg(ctx, r.pool, rec.PersonID, rec.OrgID); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE employment_records SET is_current=false, valid_to=$1, updated_at=now() WHERE person_id=$2 AND is_current=true`, rec.ValidFrom, rec.PersonID)
	if err != nil {
		return err
	}
	err = tx.QueryRow(ctx, `
		INSERT INTO employment_records (id, person_id, org_id, employee_id, job_title, job_level, employment_status, employment_type, work_arrangement, department, team, office_location, desk_number, reports_to, salary_amount, salary_currency, pay_frequency, hourly_rate, hire_date, probation_end_date, contract_start_date, contract_end_date, valid_from, is_current)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,true)
		RETURNING created_at, updated_at`,
		rec.ID, rec.PersonID, rec.OrgID, rec.EmployeeID, rec.JobTitle, rec.JobLevel, rec.EmploymentStatus, rec.EmploymentType, rec.WorkArrangement, rec.Department, rec.Team, rec.OfficeLocation, rec.DeskNumber, rec.ReportsTo, rec.SalaryAmount, rec.SalaryCurrency, rec.PayFrequency, rec.HourlyRate, rec.HireDate, rec.ProbationEndDate, rec.ContractStartDate, rec.ContractEndDate, rec.ValidFrom).Scan(&rec.CreatedAt, &rec.UpdatedAt)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *EmploymentRepository) ListByPerson(ctx context.Context, personID uuid.UUID) ([]domain.EmploymentRecord, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, person_id, org_id, employee_id, job_title, job_level, employment_status, employment_type, work_arrangement, department, team, office_location, desk_number, reports_to, salary_amount, salary_currency, pay_frequency, hourly_rate, hire_date, probation_end_date, contract_start_date, contract_end_date, valid_from, valid_to, is_current, created_at, updated_at, created_by, updated_by FROM employment_records WHERE person_id=$1 ORDER BY valid_from DESC`, personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.EmploymentRecord
	for rows.Next() {
		var e domain.EmploymentRecord
		if err := rows.Scan(&e.ID, &e.PersonID, &e.OrgID, &e.EmployeeID, &e.JobTitle, &e.JobLevel, &e.EmploymentStatus, &e.EmploymentType, &e.WorkArrangement, &e.Department, &e.Team, &e.OfficeLocation, &e.DeskNumber, &e.ReportsTo, &e.SalaryAmount, &e.SalaryCurrency, &e.PayFrequency, &e.HourlyRate, &e.HireDate, &e.ProbationEndDate, &e.ContractStartDate, &e.ContractEndDate, &e.ValidFrom, &e.ValidTo, &e.IsCurrent, &e.CreatedAt, &e.UpdatedAt, &e.CreatedBy, &e.UpdatedBy); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *EmploymentRepository) GetCurrent(ctx context.Context, personID uuid.UUID) (*domain.EmploymentRecord, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, person_id, org_id, employee_id, job_title, job_level, employment_status, employment_type, work_arrangement, department, team, office_location, desk_number, reports_to, salary_amount, salary_currency, pay_frequency, hourly_rate, hire_date, probation_end_date, contract_start_date, contract_end_date, valid_from, valid_to, is_current, created_at, updated_at, created_by, updated_by FROM employment_records WHERE person_id=$1 AND is_current=true`, personID)
	var e domain.EmploymentRecord
	err := row.Scan(&e.ID, &e.PersonID, &e.OrgID, &e.EmployeeID, &e.JobTitle, &e.JobLevel, &e.EmploymentStatus, &e.EmploymentType, &e.WorkArrangement, &e.Department, &e.Team, &e.OfficeLocation, &e.DeskNumber, &e.ReportsTo, &e.SalaryAmount, &e.SalaryCurrency, &e.PayFrequency, &e.HourlyRate, &e.HireDate, &e.ProbationEndDate, &e.ContractStartDate, &e.ContractEndDate, &e.ValidFrom, &e.ValidTo, &e.IsCurrent, &e.CreatedAt, &e.UpdatedAt, &e.CreatedBy, &e.UpdatedBy)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *EmploymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.EmploymentRecord, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, person_id, org_id, employee_id, job_title, job_level, employment_status, employment_type, work_arrangement, department, team, office_location, desk_number, reports_to, salary_amount, salary_currency, pay_frequency, hourly_rate, hire_date, probation_end_date, contract_start_date, contract_end_date, valid_from, valid_to, is_current, created_at, updated_at, created_by, updated_by FROM employment_records WHERE id=$1`, id)
	var e domain.EmploymentRecord
	err := row.Scan(&e.ID, &e.PersonID, &e.OrgID, &e.EmployeeID, &e.JobTitle, &e.JobLevel, &e.EmploymentStatus, &e.EmploymentType, &e.WorkArrangement, &e.Department, &e.Team, &e.OfficeLocation, &e.DeskNumber, &e.ReportsTo, &e.SalaryAmount, &e.SalaryCurrency, &e.PayFrequency, &e.HourlyRate, &e.HireDate, &e.ProbationEndDate, &e.ContractStartDate, &e.ContractEndDate, &e.ValidFrom, &e.ValidTo, &e.IsCurrent, &e.CreatedAt, &e.UpdatedAt, &e.CreatedBy, &e.UpdatedBy)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *EmploymentRepository) HeadcountByDept(ctx context.Context, orgID *uuid.UUID) ([]struct {
	Department string
	Count      int
}, error) {
	q := `SELECT department, COUNT(*) FROM active_employees WHERE 1=1`
	args := []any{}
	if orgID != nil {
		q += ` AND org_id=$1`
		args = append(args, *orgID)
	}
	q += ` GROUP BY department ORDER BY COUNT(*) DESC`
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []struct {
		Department string
		Count      int
	}
	for rows.Next() {
		var d string
		var c int
		_ = fmt.Sprintf
		if err := rows.Scan(&d, &c); err != nil {
			return nil, err
		}
		if d == "" {
			d = "Unassigned"
		}
		out = append(out, struct {
			Department string
			Count      int
		}{d, c})
	}
	return out, rows.Err()
}
