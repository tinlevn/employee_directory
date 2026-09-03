package repository

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"employee-directory-api/internal/domain"
	"employee-directory-api/internal/dto"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PersonRepository struct {
	pool *pgxpool.Pool
}

func NewPersonRepository(pool *pgxpool.Pool) *PersonRepository { return &PersonRepository{pool: pool} }

func (r *PersonRepository) GetOrgChart(ctx context.Context, orgID uuid.UUID) ([]domain.OrgChartNode, error) {
	q := `
		SELECT
			p.id,
			p.first_name,
			p.middle_name,
			p.last_name,
			e.job_title,
			e.department,
			p.profile_photo_url,
			e.reports_to
		FROM persons p
		JOIN employment_records e ON e.person_id = p.id AND e.is_current = true
		WHERE p.org_id = $1 AND p.is_active = true
	`
	rows, err := r.pool.Query(ctx, q, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []domain.OrgChartNode
	for rows.Next() {
		var n domain.OrgChartNode
		var fname, lname string
		var mname *string

		if err := rows.Scan(&n.ID, &fname, &mname, &lname, &n.JobTitle, &n.Department, &n.ProfilePhotoURL, &n.ReportsTo); err != nil {
			return nil, err
		}

		if mname != nil && *mname != "" {
			n.Name = fname + " " + *mname + " " + lname
		} else {
			n.Name = fname + " " + lname
		}

		nodes = append(nodes, n)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if nodes == nil {
		nodes = []domain.OrgChartNode{}
	}

	return nodes, nil
}

func (r *PersonRepository) Create(ctx context.Context, p *domain.Person) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO persons (id, org_id, first_name, middle_name, last_name, preferred_name, date_of_birth, gender, profile_photo_url, personal_email, org_email, phone_primary, address_line_1, address_line_2, city, state_province, postal_code, country, is_international, source, notes, tags)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)`,
		p.ID, p.OrgID, p.FirstName, p.MiddleName, p.LastName, p.PreferredName, p.DateOfBirth, p.Gender, p.ProfilePhotoURL, p.PersonalEmail, p.OrgEmail, p.PhonePrimary, p.AddressLine1, p.AddressLine2, p.City, p.StateProvince, p.PostalCode, p.Country, p.IsInternational, p.Source, p.Notes, p.Tags)
	return err
}

func (r *PersonRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Person, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT p.id, p.org_id, p.first_name, p.middle_name, p.last_name, p.preferred_name, p.date_of_birth, p.gender, p.profile_photo_url, p.personal_email, p.org_email, p.phone_primary, p.address_line_1, p.address_line_2, p.city, p.state_province, p.postal_code, p.country, p.is_international, p.is_active, p.archived_at, p.archive_reason, p.source, p.notes, p.tags, p.created_at, p.updated_at, p.created_by, p.updated_by,
		       e.job_title, e.department, e.team, e.office_location, e.hire_date
		FROM persons p
		LEFT JOIN employment_records e ON e.person_id = p.id AND e.is_current = true
		WHERE p.id=$1`, id)
	var p domain.Person
	err := row.Scan(&p.ID, &p.OrgID, &p.FirstName, &p.MiddleName, &p.LastName, &p.PreferredName, &p.DateOfBirth, &p.Gender, &p.ProfilePhotoURL, &p.PersonalEmail, &p.OrgEmail, &p.PhonePrimary, &p.AddressLine1, &p.AddressLine2, &p.City, &p.StateProvince, &p.PostalCode, &p.Country, &p.IsInternational, &p.IsActive, &p.ArchivedAt, &p.ArchiveReason, &p.Source, &p.Notes, &p.Tags, &p.CreatedAt, &p.UpdatedAt, &p.CreatedBy, &p.UpdatedBy, &p.CurrentJobTitle, &p.CurrentDepartment, &p.CurrentTeam, &p.CurrentOfficeLocation, &p.CurrentHireDate)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *PersonRepository) List(ctx context.Context, q dto.ListPersonsQuery) ([]domain.Person, int64, error) {
	where := []string{"1=1"}
	args := []any{}
	idx := 1

	if q.OrgID != "" {
		oid, err := uuid.Parse(q.OrgID)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid org_id: %w", err)
		}
		where = append(where, fmt.Sprintf("p.org_id = $%d", idx))
		args = append(args, oid)
		idx++
	}
	if q.Department != "" {
		where = append(where, fmt.Sprintf("EXISTS (SELECT 1 FROM employment_records e WHERE e.person_id=p.id AND e.is_current=true AND e.department ILIKE $%d)", idx))
		args = append(args, q.Department)
		idx++
	}
	if q.Search != "" {
		where = append(where, fmt.Sprintf("(p.first_name ILIKE $%d OR p.last_name ILIKE $%d OR p.preferred_name ILIKE $%d OR p.org_email ILIKE $%d OR p.personal_email ILIKE $%d OR e.job_title ILIKE $%d)", idx, idx, idx, idx, idx, idx))
		like := "%" + q.Search + "%"
		args = append(args, like)
		idx++
	}
	if q.City != "" {
		where = append(where, fmt.Sprintf("p.city ILIKE $%d", idx))
		args = append(args, q.City)
		idx++
	}
	if q.Country != "" {
		where = append(where, fmt.Sprintf("p.country ILIKE $%d", idx))
		args = append(args, q.Country)
		idx++
	}
	if q.Tag != "" {
		where = append(where, fmt.Sprintf("$%d = ANY(p.tags)", idx))
		args = append(args, q.Tag)
		idx++
	}
	if q.IsActive != nil {
		where = append(where, fmt.Sprintf("p.is_active = $%d", idx))
		args = append(args, *q.IsActive)
		idx++
	} else {
		where = append(where, "p.is_active = true")
	}
	if q.Team != "" {
		where = append(where, fmt.Sprintf("EXISTS (SELECT 1 FROM employment_records e2 WHERE e2.person_id=p.id AND e2.is_current=true AND e2.team ILIKE $%d)", idx))
		args = append(args, q.Team)
		idx++
	}
	if q.Location != "" {
		where = append(where, fmt.Sprintf("EXISTS (SELECT 1 FROM employment_records e3 WHERE e3.person_id=p.id AND e3.is_current=true AND e3.office_location ILIKE $%d)", idx))
		args = append(args, q.Location)
		idx++
	}

	whereSQL := strings.Join(where, " AND ")

	countSQL := `SELECT COUNT(*) FROM persons p LEFT JOIN employment_records e ON e.person_id=p.id AND e.is_current=true WHERE ` + whereSQL
	var total int64
	if err := r.pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Default sort is latest hire first — falls back to last_name for nulls
	orderBy := "e.hire_date DESC NULLS LAST, p.last_name ASC, p.first_name ASC"
	if q.Sort != "" {
		switch q.Sort {
		case "last_name", "-last_name", "first_name", "-first_name", "created_at", "-created_at", "city", "-city", "hire_date", "-hire_date", "hired_at", "-hired_at", "department", "-department", "job_title", "-job_title":
			dir := "ASC"
			col := q.Sort
			if strings.HasPrefix(col, "-") {
				dir = "DESC"
				col = strings.TrimPrefix(col, "-")
			}
			colMap := map[string]string{
				"last_name":  "p.last_name",
				"first_name": "p.first_name",
				"created_at": "p.created_at",
				"city":       "p.city",
				"hire_date":  "e.hire_date",
				"hired_at":   "e.hire_date",
				"department": "e.department",
				"job_title":  "e.job_title",
			}
			if qc, ok := colMap[col]; ok {
				if col == "hire_date" || col == "hired_at" {
					orderBy = qc + " " + dir + " NULLS LAST, p.last_name ASC"
				} else {
					orderBy = qc + " " + dir
				}
			}
		}
	}

	args = append(args, q.PageSize, q.Offset())
	dataSQL := fmt.Sprintf(`SELECT p.id, p.org_id, p.first_name, p.middle_name, p.last_name, p.preferred_name, p.date_of_birth, p.gender, p.profile_photo_url, p.personal_email, p.org_email, p.phone_primary, p.address_line_1, p.address_line_2, p.city, p.state_province, p.postal_code, p.country, p.is_international, p.is_active, p.archived_at, p.archive_reason, p.source, p.notes, p.tags, p.created_at, p.updated_at, p.created_by, p.updated_by,
	       e.job_title, e.department, e.team, e.office_location, e.hire_date
		FROM persons p LEFT JOIN employment_records e ON e.person_id = p.id AND e.is_current = true
		WHERE %s ORDER BY %s LIMIT $%d OFFSET $%d`, whereSQL, orderBy, idx, idx+1)

	rows, err := r.pool.Query(ctx, dataSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []domain.Person
	for rows.Next() {
		var p domain.Person
		if err := rows.Scan(&p.ID, &p.OrgID, &p.FirstName, &p.MiddleName, &p.LastName, &p.PreferredName, &p.DateOfBirth, &p.Gender, &p.ProfilePhotoURL, &p.PersonalEmail, &p.OrgEmail, &p.PhonePrimary, &p.AddressLine1, &p.AddressLine2, &p.City, &p.StateProvince, &p.PostalCode, &p.Country, &p.IsInternational, &p.IsActive, &p.ArchivedAt, &p.ArchiveReason, &p.Source, &p.Notes, &p.Tags, &p.CreatedAt, &p.UpdatedAt, &p.CreatedBy, &p.UpdatedBy, &p.CurrentJobTitle, &p.CurrentDepartment, &p.CurrentTeam, &p.CurrentOfficeLocation, &p.CurrentHireDate); err != nil {
			return nil, 0, err
		}
		out = append(out, p)
	}
	return out, total, rows.Err()
}

func (r *PersonRepository) Update(ctx context.Context, id uuid.UUID, fields map[string]any) (*domain.Person, error) {
	if len(fields) == 0 {
		return r.GetByID(ctx, id)
	}
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	set := make([]string, 0, len(keys))
	args := make([]any, 0, len(keys)+1)
	for i, k := range keys {
		set = append(set, fmt.Sprintf("%s = $%d", k, i+1))
		args = append(args, fields[k])
	}
	args = append(args, id)
	sql := fmt.Sprintf(`UPDATE persons SET %s, updated_at=now() WHERE id=$%d RETURNING id, org_id, first_name, middle_name, last_name, preferred_name, date_of_birth, gender, profile_photo_url, personal_email, org_email, phone_primary, address_line_1, address_line_2, city, state_province, postal_code, country, is_international, is_active, archived_at, archive_reason, source, notes, tags, created_at, updated_at, created_by, updated_by`, strings.Join(set, ", "), len(keys)+1)
	row := r.pool.QueryRow(ctx, sql, args...)
	var p domain.Person
	err := row.Scan(&p.ID, &p.OrgID, &p.FirstName, &p.MiddleName, &p.LastName, &p.PreferredName, &p.DateOfBirth, &p.Gender, &p.ProfilePhotoURL, &p.PersonalEmail, &p.OrgEmail, &p.PhonePrimary, &p.AddressLine1, &p.AddressLine2, &p.City, &p.StateProvince, &p.PostalCode, &p.Country, &p.IsInternational, &p.IsActive, &p.ArchivedAt, &p.ArchiveReason, &p.Source, &p.Notes, &p.Tags, &p.CreatedAt, &p.UpdatedAt, &p.CreatedBy, &p.UpdatedBy)
	if err != nil {
		return nil, err
	}
	// enrich again
	return r.GetByID(ctx, p.ID)
}

func (r *PersonRepository) SoftDelete(ctx context.Context, id uuid.UUID, reason string) (bool, error) {
	result, err := r.pool.Exec(ctx, `UPDATE persons SET is_active=false, archived_at=now(), archive_reason=$2, updated_at=now() WHERE id=$1`, id, reason)
	return result.RowsAffected() > 0, err
}
