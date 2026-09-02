package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"employee-directory-api/internal/domain"
	"employee-directory-api/internal/dto"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func integrationPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	p, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("create test pool: %v", err)
	}
	if err := p.Ping(context.Background()); err != nil {
		p.Close()
		t.Skipf("ping test database: %v", err)
	}
	t.Cleanup(p.Close)
	return p
}

func TestPersonAndEmploymentRepositories(t *testing.T) {
	pool := integrationPool(t)
	ctx := context.Background()
	orgID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	personID := uuid.New()
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM persons WHERE id=$1`, personID) })

	personRepo := NewPersonRepository(pool)
	person := &domain.Person{ID: personID, OrgID: orgID, FirstName: "Integration", LastName: "Test", IsActive: true}
	if err := personRepo.Create(ctx, person); err != nil {
		t.Fatalf("create person: %v", err)
	}

	employmentRepo := NewEmploymentRepository(pool)
	firstDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	secondDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	firstTitle := "Engineer"
	secondTitle := "Senior Engineer"
	if err := employmentRepo.CreateVersioned(ctx, &domain.EmploymentRecord{ID: uuid.New(), PersonID: personID, OrgID: orgID, JobTitle: &firstTitle, ValidFrom: firstDate}); err != nil {
		t.Fatalf("create first employment version: %v", err)
	}
	if err := employmentRepo.CreateVersioned(ctx, &domain.EmploymentRecord{ID: uuid.New(), PersonID: personID, OrgID: orgID, JobTitle: &secondTitle, ValidFrom: secondDate}); err != nil {
		t.Fatalf("create second employment version: %v", err)
	}

	current, err := employmentRepo.GetCurrent(ctx, personID)
	if err != nil {
		t.Fatalf("get current employment: %v", err)
	}
	if current == nil || current.JobTitle == nil || *current.JobTitle != secondTitle {
		t.Fatalf("expected current title %q, got %#v", secondTitle, current)
	}

	history, err := employmentRepo.ListByPerson(ctx, personID)
	if err != nil {
		t.Fatalf("list employment history: %v", err)
	}
	if len(history) != 2 {
		t.Fatalf("expected two employment versions, got %d", len(history))
	}
	if history[1].ValidTo == nil || !history[1].ValidTo.Equal(secondDate) {
		t.Fatalf("expected first version to close at %s, got %#v", secondDate, history[1].ValidTo)
	}

	persons, total, err := personRepo.List(ctx, dto.ListPersonsQuery{PaginationQuery: dto.PaginationQuery{Page: 1, PageSize: 10}, OrgID: orgID.String(), Search: "Integration"})
	if err != nil {
		t.Fatalf("list persons: %v", err)
	}
	if total != 1 || len(persons) != 1 || persons[0].CurrentJobTitle == nil {
		t.Fatalf("expected enriched person result, total=%d persons=%d", total, len(persons))
	}
}

func TestAnalyticsRepositories(t *testing.T) {
	pool := integrationPool(t)
	ctx := context.Background()
	orgID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	analytics := NewAnalyticsRepository(pool)

	headcount, err := analytics.Headcount(ctx, &orgID)
	if err != nil || len(headcount) == 0 {
		t.Fatalf("expected headcount data, err=%v rows=%d", err, len(headcount))
	}
	snapshot, err := analytics.Snapshot(ctx, "2024-12-31", &orgID)
	if err != nil || snapshot == nil || snapshot.Headcount == 0 {
		t.Fatalf("expected point-in-time snapshot, err=%v snapshot=%#v", err, snapshot)
	}
	movements, err := analytics.Movements(ctx, &orgID, "2018-01-01", "2026-12-31")
	if err != nil || len(movements) == 0 {
		t.Fatalf("expected movement data, err=%v rows=%d", err, len(movements))
	}
}

func TestAuthRepository_Integration(t *testing.T) {
	pool := integrationPool(t)
	ctx := context.Background()

	// Clean tables before test
	_, _ = pool.Exec(ctx, "TRUNCATE person_accounts, employment_records, persons, organizations CASCADE")

	orgRepo := NewOrgRepository(pool)
	personRepo := NewPersonRepository(pool)
	authRepo := NewAuthRepository(pool)

	orgID := uuid.New()
	requireNoError(t, orgRepo.Create(ctx, &domain.Organization{
		ID: orgID, Name: "Test Org", Type: "company",
	}))

	personID := uuid.New()
	requireNoError(t, personRepo.Create(ctx, &domain.Person{
		ID: personID, OrgID: orgID, FirstName: "Auth", LastName: "User", IsActive: true,
	}))

	// Create Account
	acc := &domain.PersonAccount{
		PersonID:     personID,
		Username:     "testuser",
		PasswordHash: "hashedpass",
		Role:         "admin",
		IsActive:     true,
	}
	err := authRepo.CreateAccount(ctx, acc)
	if err != nil {
		t.Fatalf("CreateAccount failed: %v", err)
	}

	if acc.ID == uuid.Nil {
		t.Errorf("expected account ID to be set")
	}

	// Get Account
	fetched, err := authRepo.GetAccountByUsername(ctx, "testuser")
	if err != nil {
		t.Fatalf("GetAccountByUsername failed: %v", err)
	}
	if fetched == nil || fetched.Username != "testuser" {
		t.Errorf("expected to fetch testuser account, got %v", fetched)
	}

	// Get Org ID
	fetchedOrg, err := authRepo.GetOrgIDByPersonID(ctx, personID)
	if err != nil {
		t.Fatalf("GetOrgIDByPersonID failed: %v", err)
	}
	if fetchedOrg != orgID {
		t.Errorf("expected org id %v, got %v", orgID, fetchedOrg)
	}
}

func TestOrgChart_Integration(t *testing.T) {
	pool := integrationPool(t)
	ctx := context.Background()

	_, _ = pool.Exec(ctx, "TRUNCATE employment_records, persons, organizations CASCADE")

	orgRepo := NewOrgRepository(pool)
	personRepo := NewPersonRepository(pool)
	empRepo := NewEmploymentRepository(pool)

	orgID := uuid.New()
	requireNoError(t, orgRepo.Create(ctx, &domain.Organization{
		ID: orgID, Name: "Chart Org", Type: "company",
	}))

	// CEO
	ceoID := uuid.New()
	requireNoError(t, personRepo.Create(ctx, &domain.Person{
		ID: ceoID, OrgID: orgID, FirstName: "Big", LastName: "Boss", IsActive: true,
	}))
	ceoTitle := "CEO"
	requireNoError(t, empRepo.CreateVersioned(ctx, &domain.EmploymentRecord{
		ID: uuid.New(), PersonID: ceoID, OrgID: orgID, IsCurrent: true,
		JobTitle: &ceoTitle, ValidFrom: time.Now(),
	}))

	// Manager
	mgrID := uuid.New()
	requireNoError(t, personRepo.Create(ctx, &domain.Person{
		ID: mgrID, OrgID: orgID, FirstName: "Middle", LastName: "Manager", IsActive: true,
	}))
	mgrTitle := "Manager"
	requireNoError(t, empRepo.CreateVersioned(ctx, &domain.EmploymentRecord{
		ID: uuid.New(), PersonID: mgrID, OrgID: orgID, IsCurrent: true,
		JobTitle: &mgrTitle, ReportsTo: &ceoID, ValidFrom: time.Now(),
	}))

	// Org Chart Query
	nodes, err := personRepo.GetOrgChart(ctx, orgID)
	if err != nil {
		t.Fatalf("GetOrgChart failed: %v", err)
	}

	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}

	foundCEO, foundMgr := false, false
	for _, n := range nodes {
		if n.ID == ceoID {
			foundCEO = true
			if n.ReportsTo != nil {
				t.Errorf("CEO should not report to anyone")
			}
		}
		if n.ID == mgrID {
			foundMgr = true
			if n.ReportsTo == nil || *n.ReportsTo != ceoID {
				t.Errorf("Manager should report to CEO")
			}
		}
	}

	if !foundCEO || !foundMgr {
		t.Errorf("missing expected nodes in org chart")
	}
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
