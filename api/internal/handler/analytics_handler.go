package handler

import (
	"employee-directory-api/internal/dto"
	"employee-directory-api/internal/middleware"
	"employee-directory-api/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AnalyticsHandler struct {
	repo      *repository.AnalyticsRepository
	validator *validator.Validate
}

func NewAnalyticsHandler(repo *repository.AnalyticsRepository, v *validator.Validate) *AnalyticsHandler {
	return &AnalyticsHandler{repo: repo, validator: v}
}

func (h *AnalyticsHandler) Headcount(c *fiber.Ctx) error {
	orgIDStr := c.Query("org_id")
	var orgID *uuid.UUID
	if orgIDStr != "" {
		if u, err := uuid.Parse(orgIDStr); err == nil {
			orgID = &u
		} else {
			return fiber.NewError(fiber.StatusBadRequest, "invalid org_id")
		}
	}
	data, err := h.repo.Headcount(c.Context(), orgID)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	if data == nil {
		data = []dto.HeadcountResponse{}
	}
	return c.JSON(data)
}

func (h *AnalyticsHandler) Attrition(c *fiber.Ctx) error {
	from := c.Query("from_date")
	to := c.Query("to_date")
	if from == "" || to == "" {
		return fiber.NewError(fiber.StatusBadRequest, "from_date and to_date required (YYYY-MM-DD)")
	}
	fromDate, fromErr := dto.ParseDate(from)
	toDate, toErr := dto.ParseDate(to)
	if fromErr != nil || toErr != nil || fromDate.After(toDate) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid date range")
	}
	orgIDStr := c.Query("org_id")
	var orgID *uuid.UUID
	if orgIDStr != "" {
		if u, err := uuid.Parse(orgIDStr); err == nil {
			orgID = &u
		} else {
			return fiber.NewError(fiber.StatusBadRequest, "invalid org_id")
		}
	}
	res, err := h.repo.Attrition(c.Context(), orgID, from, to)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	return c.JSON(res)
}

func (h *AnalyticsHandler) Movements(c *fiber.Ctx) error {
	from := c.Query("from_date")
	to := c.Query("to_date")
	if from == "" || to == "" {
		return fiber.NewError(fiber.StatusBadRequest, "from_date and to_date required")
	}
	fromDate, fromErr := dto.ParseDate(from)
	toDate, toErr := dto.ParseDate(to)
	if fromErr != nil || toErr != nil || fromDate.After(toDate) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid date range")
	}
	orgIDStr := c.Query("org_id")
	var orgID *uuid.UUID
	if orgIDStr != "" {
		if u, err := uuid.Parse(orgIDStr); err == nil {
			orgID = &u
		} else {
			return fiber.NewError(fiber.StatusBadRequest, "invalid org_id")
		}
	}
	points, err := h.repo.Movements(c.Context(), orgID, from, to)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	if points == nil {
		points = []dto.MovementPoint{}
	}
	return c.JSON(points)
}

func (h *AnalyticsHandler) Snapshot(c *fiber.Ctx) error {
	date := c.Params("date")
	if _, err := dto.ParseDate(date); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "date must use YYYY-MM-DD")
	}
	orgIDValue := c.Query("org_id")
	var orgID *uuid.UUID
	if orgIDValue != "" {
		parsed, err := uuid.Parse(orgIDValue)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid org_id")
		}
		orgID = &parsed
	}
	result, err := h.repo.Snapshot(c.Context(), date, orgID)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	return c.JSON(result)
}
