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

type OrgHandler struct {
	repo      *repository.OrgRepository
	validator *validator.Validate
}

func NewOrgHandler(repo *repository.OrgRepository, v *validator.Validate) *OrgHandler {
	return &OrgHandler{repo: repo, validator: v}
}

func (h *OrgHandler) List(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	o, err := h.repo.GetByID(c.Context(), orgID)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	list := []domain.Organization{}
	if o != nil {
		list = append(list, *o)
	}
	return c.JSON(list)
}

func (h *OrgHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	if id != middleware.GetOrgID(c) {
		return fiber.NewError(fiber.StatusForbidden, "cannot access another organization")
	}
	o, err := h.repo.GetByID(c.Context(), id)
	if err != nil {
		return err
	}
	if o == nil {
		return fiber.NewError(fiber.StatusNotFound, "organization not found")
	}
	return c.JSON(o)
}

func (h *OrgHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateOrgRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := h.validator.Struct(req); err != nil {
		return err
	}
	o := domain.Organization{ID: uuid.New(), Name: req.Name, Type: req.Type, Country: req.Country, Timezone: req.Timezone, IsActive: true}
	if err := h.repo.Create(c.Context(), &o); err != nil {
		return middleware.RepositoryError(err)
	}
	created, err := h.repo.GetByID(c.Context(), o.ID)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(created)
}
