package middleware

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"employee-directory-api/internal/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestErrorHandler_FiberError(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	app.Get("/not-found", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "custom not found message")
	})

	req := httptest.NewRequest(http.MethodGet, "/not-found", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var prob dto.ErrorResponse
	if err := json.Unmarshal(body, &prob); err != nil {
		t.Fatalf("failed to unmarshal RFC 7807 response: %v", err)
	}

	if prob.Status != 404 || prob.Detail != "custom not found message" {
		t.Fatalf("unexpected problem response: %+v", prob)
	}
}

func TestErrorHandler_InternalErrorMasking(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	app.Get("/panic-like", func(c *fiber.Ctx) error {
		return errors.New("sensitive db connection password=secret leaked")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic-like", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var prob dto.ErrorResponse
	if err := json.Unmarshal(body, &prob); err != nil {
		t.Fatalf("failed to unmarshal RFC 7807 response: %v", err)
	}

	if prob.Detail == "sensitive db connection password=secret leaked" {
		t.Fatalf("internal error detail should not leak sensitive messages, got: %s", prob.Detail)
	}
	if prob.Detail != "Internal Server Error" {
		t.Fatalf("expected 'Internal Server Error', got %q", prob.Detail)
	}
}

func TestRepositoryErrorMapping(t *testing.T) {
	// 1. pgx.ErrNoRows -> 404
	errNotFound := RepositoryError(pgx.ErrNoRows)
	var fiberErr *fiber.Error
	if !errors.As(errNotFound, &fiberErr) || fiberErr.Code != fiber.StatusNotFound {
		t.Fatalf("expected 404 for ErrNoRows, got: %v", errNotFound)
	}

	// 2. 23505 unique violation -> 409
	errUnique := RepositoryError(&pgconn.PgError{Code: "23505"})
	if !errors.As(errUnique, &fiberErr) || fiberErr.Code != fiber.StatusConflict {
		t.Fatalf("expected 409 for 23505, got: %v", errUnique)
	}

	// 3. 23503 foreign key violation -> 422
	errFK := RepositoryError(&pgconn.PgError{Code: "23503"})
	if !errors.As(errFK, &fiberErr) || fiberErr.Code != fiber.StatusUnprocessableEntity {
		t.Fatalf("expected 422 for 23503, got: %v", errFK)
	}

	// 4. 22P02 invalid input syntax -> 400
	errSyntax := RepositoryError(&pgconn.PgError{Code: "22P02"})
	if !errors.As(errSyntax, &fiberErr) || fiberErr.Code != fiber.StatusBadRequest {
		t.Fatalf("expected 400 for 22P02, got: %v", errSyntax)
	}

	// 5. unhandled generic error -> returned as-is
	customErr := errors.New("other error")
	resErr := RepositoryError(customErr)
	if resErr != customErr {
		t.Fatalf("expected generic error returned unchanged, got: %v", resErr)
	}
}

func TestNotFoundHandler(t *testing.T) {
	app := fiber.New()
	app.Use(NotFoundHandler)

	req := httptest.NewRequest(http.MethodGet, "/unknown/path", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var prob dto.ErrorResponse
	if err := json.Unmarshal(body, &prob); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if prob.Detail != "route /unknown/path not found" {
		t.Fatalf("unexpected detail: %s", prob.Detail)
	}
}
