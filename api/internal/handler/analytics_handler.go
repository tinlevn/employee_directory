package handler

import (
	"employee-directory-api/internal/dto"
	"employee-directory-api/internal/middleware"
	"employee-directory-api/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AnalyticsHandler struct {
	repo      *repository.AnalyticsRepository
	validator *validator.Validate
}

func NewAnalyticsHandler(repo *repository.AnalyticsRepository, v *validator.Validate) *AnalyticsHandler {
	return &AnalyticsHandler{repo: repo, validator: v}
}

func (h *AnalyticsHandler) Headcount(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	data, err := h.repo.Headcount(c.Context(), &orgID)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	if data == nil {
		data = []dto.HeadcountResponse{}
	}
	return c.JSON(data)
}

func (h *AnalyticsHandler) Attrition(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	from := c.Query("from")
	to := c.Query("to")
	data, err := h.repo.Attrition(c.Context(), &orgID, from, to)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	return c.JSON(data)
}

func (h *AnalyticsHandler) Movements(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	from := c.Query("from")
	to := c.Query("to")
	data, err := h.repo.Movements(c.Context(), &orgID, from, to)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	if data == nil {
		data = []dto.MovementPoint{}
	}
	return c.JSON(data)
}

func (h *AnalyticsHandler) Snapshot(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	date := c.Params("date")
	data, err := h.repo.Snapshot(c.Context(), date, &orgID)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	return c.JSON(data)
}
