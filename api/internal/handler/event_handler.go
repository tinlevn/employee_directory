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

type EventHandler struct {
	events    *repository.EventRepository
	transfers *repository.TransferRepository
	validator *validator.Validate
}

func NewEventHandler(ev *repository.EventRepository, tr *repository.TransferRepository, v *validator.Validate) *EventHandler {
	return &EventHandler{events: ev, transfers: tr, validator: v}
}

func (h *EventHandler) ListEvents(c *fiber.Ctx) error {
	pid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid person id")
	}
	var q dto.ListEventsQuery
	if err := c.QueryParser(&q); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	q.Normalize()
	if err := h.validator.Struct(q); err != nil {
		return err
	}
	events, total, err := h.events.ListByPerson(c.Context(), pid, q)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	return c.JSON(dto.PaginatedResponse[domain.StatusChangeEvent]{Data: events, Page: q.Page, PageSize: q.PageSize, Total: total, TotalPages: int((total + int64(q.PageSize) - 1) / int64(q.PageSize))})
}

func (h *EventHandler) CreateEvent(c *fiber.Ctx) error {
	pid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid person id")
	}
	var req dto.CreateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := h.validator.Struct(req); err != nil {
		return err
	}
	orgID := middleware.GetOrgID(c)
	eff, _ := time.Parse("2006-01-02", req.EffectiveDate)
	ev := domain.StatusChangeEvent{
		ID: uuid.New(), PersonID: pid, OrgID: orgID,
		EventType: req.EventType, Context: "employment",
		FromStatus: req.FromStatus, ToStatus: req.ToStatus,
		FromDepartment: req.FromDepartment, ToDepartment: req.ToDepartment,
		FromTitle: req.FromTitle, ToTitle: req.ToTitle,
		FromLocation: req.FromLocation, ToLocation: req.ToLocation,
		Reason: req.Reason, ReasonCode: req.ReasonCode, IsVoluntary: req.IsVoluntary,
		EffectiveDate: eff, DocumentURLs: req.DocumentURLs, Notes: req.Notes,
	}
	if req.Context != nil {
		ev.Context = *req.Context
	}
	if req.InitiatedBy != nil {
		if u, err := uuid.Parse(*req.InitiatedBy); err == nil {
			ev.InitiatedBy = &u
		}
	}
	if req.ApprovedBy != nil {
		if u, err := uuid.Parse(*req.ApprovedBy); err == nil {
			ev.ApprovedBy = &u
		}
	}
	if req.WitnessedBy != nil {
		if u, err := uuid.Parse(*req.WitnessedBy); err == nil {
			ev.WitnessedBy = &u
		}
	}
	if err := h.events.Create(c.Context(), &ev); err != nil {
		if errors.Is(err, repository.ErrOrganizationMismatch) {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "event references a different organization")
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusNotFound, "person or referenced actor not found")
		}
		return middleware.RepositoryError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(ev)
}

func (h *EventHandler) ListTransfers(c *fiber.Ctx) error {
	pid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid person id")
	}
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	list, total, err := h.transfers.ListByPerson(c.Context(), pid, page, pageSize)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	return c.JSON(dto.PaginatedResponse[domain.TransferRecord]{Data: list, Page: page, PageSize: pageSize, Total: total, TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize))})
}

func (h *EventHandler) CreateTransfer(c *fiber.Ctx) error {
	pid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid person id")
	}
	var req dto.CreateTransferRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := h.validator.Struct(req); err != nil {
		return err
	}
	orgID := middleware.GetOrgID(c)
	eff, _ := time.Parse("2006-01-02", req.EffectiveDate)
	t := domain.TransferRecord{
		ID: uuid.New(), PersonID: pid, OrgID: orgID,
		TransferType: req.TransferType, FromDepartment: req.FromDepartment, ToDepartment: req.ToDepartment,
		FromLocation: req.FromLocation, ToLocation: req.ToLocation,
		FromTitle: req.FromTitle, ToTitle: req.ToTitle,
		EffectiveDate: eff, Reason: req.Reason, Notes: req.Notes,
	}
	if req.FromManagerID != nil {
		if u, err := uuid.Parse(*req.FromManagerID); err == nil {
			t.FromManagerID = &u
		}
	}
	if req.ToManagerID != nil {
		if u, err := uuid.Parse(*req.ToManagerID); err == nil {
			t.ToManagerID = &u
		}
	}
	if req.ApprovedBy != nil {
		if u, err := uuid.Parse(*req.ApprovedBy); err == nil {
			t.ApprovedBy = &u
		}
	}
	if err := h.transfers.Create(c.Context(), &t); err != nil {
		if errors.Is(err, repository.ErrOrganizationMismatch) {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "transfer references a different organization")
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusNotFound, "person or referenced manager not found")
		}
		return middleware.RepositoryError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(t)
}
