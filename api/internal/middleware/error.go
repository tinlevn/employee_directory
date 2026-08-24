package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"employee-directory-api/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "internal server error"

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		msg = e.Message
	} else if ve, ok := err.(validator.ValidationErrors); ok {
		code = fiber.StatusBadRequest
		errs := make(map[string]string)
		for _, fe := range ve {
			errs[fe.Field()] = fmt.Sprintf("failed on '%s'", fe.Tag())
		}
		return c.Status(code).JSON(dto.ErrorResponse{
			Type:    "https://httpstatuses.com/400",
			Title:   "Validation failed",
			Status:  code,
			Detail:  "One or more fields failed validation",
			Errors:  errs,
			TraceID: c.GetRespHeader("X-Request-Id"),
		})
	} else if err != nil {
		msg = http.StatusText(code)
	}
	if code >= fiber.StatusInternalServerError && code != fiber.StatusNotImplemented {
		msg = http.StatusText(code)
	}

	return c.Status(code).JSON(dto.ErrorResponse{
		Type:    fmt.Sprintf("https://httpstatuses.com/%d", code),
		Title:   http.StatusText(code),
		Status:  code,
		Detail:  msg,
		TraceID: c.GetRespHeader("X-Request-Id"),
	})
}

// RepositoryError maps stable PostgreSQL conditions to public HTTP errors without exposing SQL details.
func RepositoryError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return fiber.NewError(fiber.StatusNotFound, "resource not found")
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return fiber.NewError(fiber.StatusConflict, "resource already exists")
		case "23503", "23514":
			return fiber.NewError(fiber.StatusUnprocessableEntity, "resource violates a data constraint")
		case "22P02":
			return fiber.NewError(fiber.StatusBadRequest, "invalid value")
		}
	}
	return err
}

func NotFoundHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
		Type:   "https://httpstatuses.com/404",
		Title:  "Not Found",
		Status: 404,
		Detail: fmt.Sprintf("route %s not found", c.Path()),
	})
}
