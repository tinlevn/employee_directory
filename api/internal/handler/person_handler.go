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

type PersonHandler struct {
	repo      *repository.PersonRepository
	validator *validator.Validate
}

func NewPersonHandler(repo *repository.PersonRepository, v *validator.Validate) *PersonHandler {
	return &PersonHandler{repo: repo, validator: v}
}

func (h *PersonHandler) List(c *fiber.Ctx) error {
	var q dto.ListPersonsQuery
	if err := c.QueryParser(&q); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	q.Normalize()
	if err := h.validator.Struct(q); err != nil {
		return err
	}
	q.OrgID = middleware.GetOrgID(c).String()
	persons, total, err := h.repo.List(c.Context(), q)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	return c.JSON(dto.PaginatedResponse[domain.Person]{
		Data: persons, Page: q.Page, PageSize: q.PageSize, Total: total, TotalPages: int((total + int64(q.PageSize) - 1) / int64(q.PageSize)),
	})
}

func (h *PersonHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	p, err := h.repo.GetByID(c.Context(), id)
	if err != nil {
		return err
	}
	if p == nil {
		return fiber.NewError(fiber.StatusNotFound, "person not found")
	}
	return c.JSON(p)
}

func (h *PersonHandler) Create(c *fiber.Ctx) error {
	var req dto.CreatePersonRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := h.validator.Struct(req); err != nil {
		return err
	}
	orgID := middleware.GetOrgID(c)
	p := domain.Person{
		ID: uuid.New(), OrgID: orgID,
		FirstName: req.FirstName, MiddleName: req.MiddleName, LastName: req.LastName,
		PreferredName: req.PreferredName, Gender: req.Gender, ProfilePhotoURL: req.ProfilePhotoURL,
		PersonalEmail: req.PersonalEmail, OrgEmail: req.OrgEmail, PhonePrimary: req.PhonePrimary,
		AddressLine1: req.AddressLine1, AddressLine2: req.AddressLine2, City: req.City, StateProvince: req.StateProvince, PostalCode: req.PostalCode, Country: req.Country,
		Source: req.Source, Notes: req.Notes, Tags: req.Tags, IsActive: true,
	}
	if req.IsInternational != nil {
		p.IsInternational = *req.IsInternational
	}
	if req.DateOfBirth != nil {
		if t, err := time.Parse("2006-01-02", *req.DateOfBirth); err == nil {
			p.DateOfBirth = &t
		}
	}
	if err := h.repo.Create(c.Context(), &p); err != nil {
		return middleware.RepositoryError(err)
	}
	created, err := h.repo.GetByID(c.Context(), p.ID)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h *PersonHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	var req dto.UpdatePersonRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := h.validator.Struct(req); err != nil {
		return err
	}
	fields := map[string]any{}
	if req.FirstName != nil {
		fields["first_name"] = *req.FirstName
	}
	if req.MiddleName != nil {
		fields["middle_name"] = *req.MiddleName
	}
	if req.LastName != nil {
		fields["last_name"] = *req.LastName
	}
	if req.PreferredName != nil {
		fields["preferred_name"] = *req.PreferredName
	}
	if req.DateOfBirth != nil {
		if t, err := time.Parse("2006-01-02", *req.DateOfBirth); err == nil {
			fields["date_of_birth"] = t
		}
	}
	if req.Gender != nil {
		fields["gender"] = *req.Gender
	}
	if req.ProfilePhotoURL != nil {
		fields["profile_photo_url"] = *req.ProfilePhotoURL
	}
	if req.PersonalEmail != nil {
		fields["personal_email"] = *req.PersonalEmail
	}
	if req.OrgEmail != nil {
		fields["org_email"] = *req.OrgEmail
	}
	if req.PhonePrimary != nil {
		fields["phone_primary"] = *req.PhonePrimary
	}
	if req.AddressLine1 != nil {
		fields["address_line_1"] = *req.AddressLine1
	}
	if req.AddressLine2 != nil {
		fields["address_line_2"] = *req.AddressLine2
	}
	if req.City != nil {
		fields["city"] = *req.City
	}
	if req.StateProvince != nil {
		fields["state_province"] = *req.StateProvince
	}
	if req.PostalCode != nil {
		fields["postal_code"] = *req.PostalCode
	}
	if req.Country != nil {
		fields["country"] = *req.Country
	}
	if req.IsInternational != nil {
		fields["is_international"] = *req.IsInternational
	}
	if req.Source != nil {
		fields["source"] = *req.Source
	}
	if req.Notes != nil {
		fields["notes"] = *req.Notes
	}
	if req.Tags != nil {
		fields["tags"] = *req.Tags
	}

	updated, err := h.repo.Update(c.Context(), id, fields)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusNotFound, "person not found")
		}
		return middleware.RepositoryError(err)
	}
	return c.JSON(updated)
}

func (h *PersonHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	reason := c.Query("reason", "deactivated")
	found, err := h.repo.SoftDelete(c.Context(), id, reason)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "person not found")
	}
	return c.JSON(dto.MessageResponse{Message: "person deactivated", ID: id.String()})
}
