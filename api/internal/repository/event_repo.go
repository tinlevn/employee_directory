package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"employee-directory-api/internal/domain"
	"employee-directory-api/internal/dto"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct {
	pool *pgxpool.Pool
}

func NewEventRepository(pool *pgxpool.Pool) *EventRepository { return &EventRepository{pool: pool} }

func (r *EventRepository) Create(ctx context.Context, e *domain.StatusChangeEvent) error {
	if err := ensurePersonInOrg(ctx, r.pool, e.PersonID, e.OrgID); err != nil {
		return err
	}
	for _, actor := range []*uuid.UUID{e.InitiatedBy, e.ApprovedBy, e.WitnessedBy} {
		if actor != nil {
			if err := ensurePersonInOrg(ctx, r.pool, *actor, e.OrgID); err != nil {
				return err
			}
		}
	}
	return r.pool.QueryRow(ctx, `
		INSERT INTO status_change_events (id, person_id, org_id, event_type, context, from_status, to_status, from_department, to_department, from_title, to_title, from_location, to_location, reason, reason_code, is_voluntary, effective_date, initiated_by, approved_by, witnessed_by, linked_record_id, linked_record_type, document_urls, notes)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24)
		RETURNING recorded_at, created_at`,
		e.ID, e.PersonID, e.OrgID, e.EventType, e.Context, e.FromStatus, e.ToStatus, e.FromDepartment, e.ToDepartment, e.FromTitle, e.ToTitle, e.FromLocation, e.ToLocation, e.Reason, e.ReasonCode, e.IsVoluntary, e.EffectiveDate, e.InitiatedBy, e.ApprovedBy, e.WitnessedBy, e.LinkedRecordID, e.LinkedRecordType, e.DocumentURLs, e.Notes).Scan(&e.RecordedAt, &e.CreatedAt)
}

func (r *EventRepository) ListByPerson(ctx context.Context, personID uuid.UUID, q dto.ListEventsQuery) ([]domain.StatusChangeEvent, int64, error) {
	where := []string{"person_id = $1"}
	args := []any{personID}
	idx := 2
	if q.EventType != "" {
		where = append(where, fmt.Sprintf("event_type = $%d", idx))
		args = append(args, q.EventType)
		idx++
	}
	if q.FromDate != nil && *q.FromDate != "" {
		where = append(where, fmt.Sprintf("effective_date >= $%d", idx))
		t, _ := dto.ParseDate(*q.FromDate)
		args = append(args, t)
		idx++
	}
	if q.ToDate != nil && *q.ToDate != "" {
		where = append(where, fmt.Sprintf("effective_date <= $%d", idx))
		t, _ := dto.ParseDate(*q.ToDate)
		args = append(args, t)
		idx++
	}
	countSQL := `SELECT COUNT(*) FROM status_change_events WHERE ` + strings.Join(where, " AND ")
	var total int64
	if err := r.pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	dataSQL := fmt.Sprintf(`SELECT id, person_id, org_id, event_type, context, from_status, to_status, from_department, to_department, from_title, to_title, from_location, to_location, reason, reason_code, is_voluntary, effective_date, recorded_at, initiated_by, approved_by, witnessed_by, linked_record_id, linked_record_type, document_urls, notes, created_at, created_by FROM status_change_events WHERE %s ORDER BY effective_date DESC, recorded_at DESC LIMIT $%d OFFSET $%d`, strings.Join(where, " AND "), idx, idx+1)
	args = append(args, q.PageSize, q.Offset())
	rows, err := r.pool.Query(ctx, dataSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []domain.StatusChangeEvent
	for rows.Next() {
		var e domain.StatusChangeEvent
		if err := rows.Scan(&e.ID, &e.PersonID, &e.OrgID, &e.EventType, &e.Context, &e.FromStatus, &e.ToStatus, &e.FromDepartment, &e.ToDepartment, &e.FromTitle, &e.ToTitle, &e.FromLocation, &e.ToLocation, &e.Reason, &e.ReasonCode, &e.IsVoluntary, &e.EffectiveDate, &e.RecordedAt, &e.InitiatedBy, &e.ApprovedBy, &e.WitnessedBy, &e.LinkedRecordID, &e.LinkedRecordType, &e.DocumentURLs, &e.Notes, &e.CreatedAt, &e.CreatedBy); err != nil {
			return nil, 0, err
		}
		out = append(out, e)
	}
	return out, total, rows.Err()
}

type TransferRepository struct {
	pool *pgxpool.Pool
}

func NewTransferRepository(pool *pgxpool.Pool) *TransferRepository {
	return &TransferRepository{pool: pool}
}

func (r *TransferRepository) Create(ctx context.Context, t *domain.TransferRecord) error {
	if err := ensurePersonInOrg(ctx, r.pool, t.PersonID, t.OrgID); err != nil {
		return err
	}
	for _, manager := range []*uuid.UUID{t.FromManagerID, t.ToManagerID, t.ApprovedBy} {
		if manager != nil {
			if err := ensurePersonInOrg(ctx, r.pool, *manager, t.OrgID); err != nil {
				return err
			}
		}
	}
	return r.pool.QueryRow(ctx, `INSERT INTO transfer_records (id, person_id, org_id, transfer_type, from_department, to_department, from_location, to_location, from_manager_id, to_manager_id, from_title, to_title, effective_date, reason, approved_by, notes) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) RETURNING created_at`,
		t.ID, t.PersonID, t.OrgID, t.TransferType, t.FromDepartment, t.ToDepartment, t.FromLocation, t.ToLocation, t.FromManagerID, t.ToManagerID, t.FromTitle, t.ToTitle, t.EffectiveDate, t.Reason, t.ApprovedBy, t.Notes).Scan(&t.CreatedAt)
}

func (r *TransferRepository) ListByPerson(ctx context.Context, personID uuid.UUID, page, pageSize int) ([]domain.TransferRecord, int64, error) {
	var total int64
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM transfer_records WHERE person_id=$1`, personID).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, `SELECT id, person_id, org_id, transfer_type, from_department, to_department, from_location, to_location, from_manager_id, to_manager_id, from_title, to_title, effective_date, reason, approved_by, notes, created_at, created_by FROM transfer_records WHERE person_id=$1 ORDER BY effective_date DESC LIMIT $2 OFFSET $3`, personID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []domain.TransferRecord
	for rows.Next() {
		var t domain.TransferRecord
		if err := rows.Scan(&t.ID, &t.PersonID, &t.OrgID, &t.TransferType, &t.FromDepartment, &t.ToDepartment, &t.FromLocation, &t.ToLocation, &t.FromManagerID, &t.ToManagerID, &t.FromTitle, &t.ToTitle, &t.EffectiveDate, &t.Reason, &t.ApprovedBy, &t.Notes, &t.CreatedAt, &t.CreatedBy); err != nil {
			return nil, 0, err
		}
		out = append(out, t)
	}
	return out, total, rows.Err()
}

// Analytics repository
type AnalyticsRepository struct {
	pool *pgxpool.Pool
}

func NewAnalyticsRepository(pool *pgxpool.Pool) *AnalyticsRepository {
	return &AnalyticsRepository{pool: pool}
}

func (r *AnalyticsRepository) Snapshot(ctx context.Context, date string, orgID *uuid.UUID) (*dto.SnapshotResponse, error) {
	day, err := dto.ParseDate(date)
	if err != nil {
		return nil, err
	}

	args := []any{day}
	orgFilter := ""
	if orgID != nil {
		orgFilter = " AND p.org_id=$2"
		args = append(args, *orgID)
	}
	base := `FROM persons p JOIN employment_records e ON e.person_id=p.id AND e.org_id=p.org_id
		WHERE p.is_active=true AND e.valid_from <= $1 AND (e.valid_to IS NULL OR e.valid_to > $1)` + orgFilter

	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) `+base, args...).Scan(&total); err != nil {
		return nil, err
	}
	departmentCounts, err := r.snapshotGroup(ctx, base, args, "e.department")
	if err != nil {
		return nil, err
	}
	byDepartment := make([]dto.HeadcountResponse, 0, len(departmentCounts))
	for _, row := range departmentCounts {
		byDepartment = append(byDepartment, dto.HeadcountResponse{Department: row.Name, Count: row.Count})
	}
	byLocation, err := r.snapshotGroup(ctx, base, args, "e.office_location")
	if err != nil {
		return nil, err
	}
	resultOrg := ""
	if orgID != nil {
		resultOrg = orgID.String()
	}
	return &dto.SnapshotResponse{Date: date, OrgID: resultOrg, Headcount: total, ByDept: byDepartment, ByLocation: byLocation}, nil
}

func (r *AnalyticsRepository) snapshotGroup(ctx context.Context, base string, args []any, column string) ([]dto.DimensionCount, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+column+`, COUNT(*) `+base+` GROUP BY `+column+` ORDER BY COUNT(*) DESC`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []dto.DimensionCount
	for rows.Next() {
		var name *string
		var count int
		if err := rows.Scan(&name, &count); err != nil {
			return nil, err
		}
		label := "Unassigned"
		if name != nil && *name != "" {
			label = *name
		}
		result = append(result, dto.DimensionCount{Name: label, Count: count})
	}
	return result, rows.Err()
}

func (r *AnalyticsRepository) Headcount(ctx context.Context, orgID *uuid.UUID) ([]dto.HeadcountResponse, error) {
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
	var out []dto.HeadcountResponse
	for rows.Next() {
		var d *string
		var c int
		if err := rows.Scan(&d, &c); err != nil {
			return nil, err
		}
		dept := "Unassigned"
		if d != nil && *d != "" {
			dept = *d
		}
		out = append(out, dto.HeadcountResponse{Department: dept, Count: c})
	}
	return out, rows.Err()
}

func (r *AnalyticsRepository) Attrition(ctx context.Context, orgID *uuid.UUID, from, to string) (*dto.AttritionResponse, error) {
	fd, _ := dto.ParseDate(from)
	td, _ := dto.ParseDate(to)
	// total exits
	var total, voluntary, involuntary int64
	var args []any
	qBase := `FROM status_change_events WHERE event_type IN ('RESIGNED','TERMINATED','LAID_OFF') AND effective_date >= $1 AND effective_date <= $2`
	args = append(args, fd, td)
	if orgID != nil {
		qBase += ` AND org_id=$3`
		args = append(args, *orgID)
	}
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) `+qBase, args...).Scan(&total); err != nil {
		return nil, err
	}
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) `+qBase+` AND is_voluntary=true`, args...).Scan(&voluntary); err != nil {
		return nil, err
	}
	involuntary = total - voluntary

	// avg headcount from snapshots
	var avg *float64
	avgQ := `SELECT AVG(active_count) FROM headcount_snapshots WHERE snapshot_date >= $1 AND snapshot_date <= $2`
	avgArgs := []any{fd, td}
	if orgID != nil {
		avgQ += ` AND org_id=$3`
		avgArgs = append(avgArgs, *orgID)
	}
	_ = r.pool.QueryRow(ctx, avgQ, avgArgs...).Scan(&avg)
	avgVal := 0.0
	if avg != nil {
		avgVal = *avg
	}
	rate := 0.0
	if avgVal > 0 {
		rate = float64(total) / avgVal * 100
	}
	return &dto.AttritionResponse{
		PeriodStart: from, PeriodEnd: to,
		TotalExits: total, VoluntaryExits: voluntary, InvoluntaryExits: involuntary,
		AvgHeadcount: avgVal, AttritionRatePct: rate,
	}, nil
}

func (r *AnalyticsRepository) Movements(ctx context.Context, orgID *uuid.UUID, from, to string) ([]dto.MovementPoint, error) {
	fd, _ := dto.ParseDate(from)
	td, _ := dto.ParseDate(to)
	// group by date using status_change_events
	q := `SELECT effective_date::date as d,
		COUNT(*) FILTER (WHERE event_type='HIRED') as hires,
		COUNT(*) FILTER (WHERE event_type IN ('RESIGNED','TERMINATED','LAID_OFF')) as exits
		FROM status_change_events WHERE effective_date >= $1 AND effective_date <= $2`
	args := []any{fd, td}
	if orgID != nil {
		q += ` AND org_id=$3`
		args = append(args, *orgID)
	}
	q += ` GROUP BY d ORDER BY d`
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []dto.MovementPoint
	for rows.Next() {
		var d time.Time
		var hires, exits int
		if err := rows.Scan(&d, &hires, &exits); err != nil {
			return nil, err
		}
		out = append(out, dto.MovementPoint{Date: d.Format("2006-01-02"), NewEntries: hires, Exits: exits, NetChange: hires - exits})
	}
	return out, rows.Err()
}
