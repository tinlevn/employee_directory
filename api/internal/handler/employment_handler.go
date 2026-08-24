package handler

import (
	"errors"
	"time"

	"employee-directory-api/internal/domain"
	"employee-directory-api/internal/dto"
	"employee-directory-api/internal/middleware"
	"employee-directory-api/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type EmploymentHandler struct {
	repo      *repository.EmploymentRepository
	validator *validator.Validate
}

func NewEmploymentHandler(repo *repository.EmploymentRepository, v *validator.Validate) *EmploymentHandler {
	return &EmploymentHandler{repo: repo, validator: v}
}

func (h *EmploymentHandler) List(c *fiber.Ctx) error {
	pid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid person id")
	}
	list, err := h.repo.ListByPerson(c.Context(), pid)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	if list == nil {
		list = []domain.EmploymentRecord{}
	}
	return c.JSON(list)
}

func (h *EmploymentHandler) GetCurrent(c *fiber.Ctx) error {
	pid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid person id")
	}
	cur, err := h.repo.GetCurrent(c.Context(), pid)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	if cur == nil {
		return fiber.NewError(fiber.StatusNotFound, "no current employment record")
	}
	return c.JSON(cur)
}

func (h *EmploymentHandler) Create(c *fiber.Ctx) error {
	pid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid person id")
	}
	var req dto.CreateEmploymentRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := h.validator.Struct(req); err != nil {
		return err
	}
	orgID, _ := uuid.Parse(req.OrgID)
	validFrom, _ := time.Parse("2006-01-02", req.ValidFrom)
	rec := domain.EmploymentRecord{
		ID: uuid.New(), PersonID: pid, OrgID: orgID,
		EmployeeID: req.EmployeeID, JobTitle: req.JobTitle, JobLevel: req.JobLevel,
		EmploymentStatus: req.EmploymentStatus, EmploymentType: req.EmploymentType, WorkArrangement: req.WorkArrangement,
		Department: req.Department, Team: req.Team, OfficeLocation: req.OfficeLocation, DeskNumber: req.DeskNumber,
		SalaryAmount: req.SalaryAmount, SalaryCurrency: req.SalaryCurrency, PayFrequency: req.PayFrequency, HourlyRate: req.HourlyRate,
		ValidFrom: validFrom, IsCurrent: true,
	}
	if req.ReportsTo != nil {
		if u, err := uuid.Parse(*req.ReportsTo); err == nil {
			rec.ReportsTo = &u
		}
	}
	if req.HireDate != nil {
		if t, err := time.Parse("2006-01-02", *req.HireDate); err == nil {
			rec.HireDate = &t
		}
	}
	if req.ProbationEndDate != nil {
		if t, err := time.Parse("2006-01-02", *req.ProbationEndDate); err == nil {
			rec.ProbationEndDate = &t
		}
	}
	if req.ContractStartDate != nil {
		if t, err := time.Parse("2006-01-02", *req.ContractStartDate); err == nil {
			rec.ContractStartDate = &t
		}
	}
	if req.ContractEndDate != nil {
		if t, err := time.Parse("2006-01-02", *req.ContractEndDate); err == nil {
			rec.ContractEndDate = &t
		}
	}
	if err := h.repo.CreateVersioned(c.Context(), &rec); err != nil {
		if errors.Is(err, repository.ErrOrganizationMismatch) {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "person and employment organization must match")
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusNotFound, "person not found")
		}
		return middleware.RepositoryError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(rec)
}
