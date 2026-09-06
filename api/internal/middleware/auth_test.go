package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"employee-directory-api/internal/auth"
	"employee-directory-api/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func setupAuthTestApp(svc *auth.Service) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	app.Get("/protected", RequireAuth(svc), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"org_id":     GetOrgID(c),
			"account_id": GetAccountID(c),
			"person_id":  GetPersonID(c),
			"role":       GetRole(c),
		})
	})

	app.Get("/admin-only", RequireAuth(svc), RequireRole("admin"), func(c *fiber.Ctx) error {
		return c.SendString("admin granted")
	})

	app.Get("/staff", RequireAuth(svc), RequireRole("admin", "manager"), func(c *fiber.Ctx) error {
		return c.SendString("staff granted")
	})

	return app
}

func TestRequireAuth_MissingHeader(t *testing.T) {
	svc, _ := auth.NewService("test-secret", "1h")
	app := setupAuthTestApp(svc)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestRequireAuth_InvalidHeaderFormat(t *testing.T) {
	svc, _ := auth.NewService("test-secret", "1h")
	app := setupAuthTestApp(svc)

	cases := []string{
		"Basic xyz",
		"Bearer",
		"Bearer ",
		"Token abc",
	}

	for _, header := range cases {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", header)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("app.Test failed: %v", err)
		}
		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("for header %q, expected 401, got %d", header, resp.StatusCode)
		}
	}
}

func TestRequireAuth_InvalidToken(t *testing.T) {
	svc, _ := auth.NewService("test-secret", "1h")
	app := setupAuthTestApp(svc)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer bad.token.here")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestRequireAuth_ValidToken(t *testing.T) {
	svc, _ := auth.NewService("test-secret", "1h")
	app := setupAuthTestApp(svc)

	accID := uuid.New()
	personID := uuid.New()
	orgID := uuid.New()
	acc := &domain.PersonAccount{
		ID:       accID,
		Username: "alice",
		PersonID: personID,
		Role:     "employee",
	}

	tok, err := svc.IssueToken(acc, orgID)
	if err != nil {
		t.Fatalf("IssueToken failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tok))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRequireRole_AuthorizedAndForbidden(t *testing.T) {
	svc, _ := auth.NewService("test-secret", "1h")
	app := setupAuthTestApp(svc)

	managerAcc := &domain.PersonAccount{
		ID:       uuid.New(),
		Username: "manager1",
		PersonID: uuid.New(),
		Role:     "manager",
	}
	adminAcc := &domain.PersonAccount{
		ID:       uuid.New(),
		Username: "admin1",
		PersonID: uuid.New(),
		Role:     "admin",
	}

	managerToken, _ := svc.IssueToken(managerAcc, uuid.New())
	adminToken, _ := svc.IssueToken(adminAcc, uuid.New())

	// Manager accessing /admin-only -> 403
	reqAdminOnly := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	reqAdminOnly.Header.Set("Authorization", "Bearer "+managerToken)
	resp, err := app.Test(reqAdminOnly)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for manager on admin route, got %d", resp.StatusCode)
	}

	// Admin accessing /admin-only -> 200
	reqAdminOK := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	reqAdminOK.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err = app.Test(reqAdminOK)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK for admin on admin route, got %d", resp.StatusCode)
	}

	// Manager accessing /staff (accepts admin, manager) -> 200
	reqStaffOK := httptest.NewRequest(http.MethodGet, "/staff", nil)
	reqStaffOK.Header.Set("Authorization", "Bearer "+managerToken)
	resp, err = app.Test(reqStaffOK)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK for manager on staff route, got %d", resp.StatusCode)
	}
}

func TestGetAccountID_InvalidSubjectNoPanic(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		// Set claims with invalid UUID subject
		c.Locals(claimsKey, &auth.Claims{
			Role: "admin",
		})
		// Subject is empty string (invalid UUID)
		accID := GetAccountID(c)
		if accID != uuid.Nil {
			t.Fatalf("expected uuid.Nil, got %v", accID)
		}
		return c.SendStatus(200)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
