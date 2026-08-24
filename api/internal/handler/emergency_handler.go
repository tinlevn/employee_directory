package handler

import (
	"employee-directory-api/internal/domain"
	"employee-directory-api/internal/dto"
	"employee-directory-api/internal/middleware"
	"employee-directory-api/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type EmergencyHandler struct {
	repo      *repository.EmergencyContactRepository
	validator *validator.Validate
}

func NewEmergencyHandler(repo *repository.EmergencyContactRepository, v *validator.Validate) *EmergencyHandler {
	return &EmergencyHandler{repo: repo, validator: v}
}

func (h *EmergencyHandler) Get(c *fiber.Ctx) error {
	pid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid person id")
	}
	ec, err := h.repo.GetByPerson(c.Context(), pid)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	if ec == nil {
		return fiber.NewError(fiber.StatusNotFound, "no emergency contact")
	}
	return c.JSON(ec)
}

func (h *EmergencyHandler) Upsert(c *fiber.Ctx) error {
	pid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid person id")
	}
	var req dto.CreateEmergencyContactRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := h.validator.Struct(req); err != nil {
		return err
	}
	// if exists, keep same id
	existing, err := h.repo.GetByPerson(c.Context(), pid)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	id := uuid.New()
	if existing != nil {
		id = existing.ID
	}
	ec := domain.EmergencyContact{ID: id, PersonID: pid, Name: req.Name, Phone: req.Phone, Email: req.Email, Relationship: req.Relationship}
	saved, err := h.repo.Upsert(c.Context(), &ec)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(saved)
}

func (h *EmergencyHandler) Update(c *fiber.Ctx) error {
	pid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid person id")
	}
	existing, err := h.repo.GetByPerson(c.Context(), pid)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	if existing == nil {
		return fiber.NewError(fiber.StatusNotFound, "no emergency contact to update")
	}
	var req dto.UpdateEmergencyContactRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := h.validator.Struct(req); err != nil {
		return err
	}
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Phone != nil {
		existing.Phone = req.Phone
	}
	if req.Email != nil {
		existing.Email = req.Email
	}
	if req.Relationship != nil {
		existing.Relationship = req.Relationship
	}
	saved, err := h.repo.Upsert(c.Context(), existing)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	return c.JSON(saved)
}

func (h *EmergencyHandler) Delete(c *fiber.Ctx) error {
	pid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid person id")
	}
	found, err := h.repo.Delete(c.Context(), pid)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "no emergency contact")
	}
	return c.JSON(dto.MessageResponse{Message: "emergency contact deleted"})
}
