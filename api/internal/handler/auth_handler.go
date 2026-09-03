package handler

import (
	"errors"

	"employee-directory-api/internal/auth"
	"employee-directory-api/internal/domain"
	"employee-directory-api/internal/middleware"
	"employee-directory-api/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AuthHandler struct {
	repo      *repository.AuthRepository
	svc       *auth.Service
	validator *validator.Validate
}

func NewAuthHandler(repo *repository.AuthRepository, svc *auth.Service, v *validator.Validate) *AuthHandler {
	return &AuthHandler{repo: repo, svc: svc, validator: v}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := h.validator.Struct(req); err != nil {
		return err
	}

	role := req.Role
	if role == "" {
		role = "staff"
	}
	if role == "admin" || role == "manager" {
		return fiber.NewError(fiber.StatusForbidden, "elevated roles must be provisioned by an administrator")
	}

	personID, err := uuid.Parse(req.PersonID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid person_id")
	}
	orgID, err := h.repo.GetOrgIDByPersonID(c.Context(), personID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusNotFound, "person not found")
		}
		return middleware.RepositoryError(err)
	}

	if existing, err := h.repo.GetAccountByUsername(c.Context(), req.Username); err != nil {
		return middleware.RepositoryError(err)
	} else if existing != nil {
		return fiber.NewError(fiber.StatusConflict, "username already taken")
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return err
	}

	account := &domain.PersonAccount{
		PersonID:     personID,
		Username:     req.Username,
		PasswordHash: hash,
		Role:         role,
		IsActive:     true,
	}
	if err := h.repo.CreateAccount(c.Context(), account); err != nil {
		return middleware.RepositoryError(err)
	}

	token, err := h.svc.IssueToken(account, orgID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(domain.AuthResponse{Token: token, Account: account, OrgID: orgID})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := h.validator.Struct(req); err != nil {
		return err
	}

	account, err := h.repo.GetAccountByUsername(c.Context(), req.Username)
	if err != nil {
		return middleware.RepositoryError(err)
	}
	if account == nil || !account.IsActive || account.AccountLocked || !auth.CheckPassword(account.PasswordHash, req.Password) {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}

	orgID, err := h.repo.GetOrgIDByPersonID(c.Context(), account.PersonID)
	if err != nil {
		return middleware.RepositoryError(err)
	}

	token, err := h.svc.IssueToken(account, orgID)
	if err != nil {
		return err
	}
	_ = h.repo.UpdateLastLogin(c.Context(), account.ID)
	return c.JSON(domain.AuthResponse{Token: token, Account: account, OrgID: orgID})
}
