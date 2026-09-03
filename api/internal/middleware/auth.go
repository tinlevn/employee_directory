package middleware

import (
	"strings"

	"employee-directory-api/internal/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const claimsKey = "auth_claims"

// RequireAuth validates the Bearer token and stores the parsed claims in locals.
func RequireAuth(svc *auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get(fiber.HeaderAuthorization)
		if header == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid authorization header")
		}
		claims, err := svc.ParseToken(strings.TrimSpace(parts[1]))
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired token")
		}
		c.Locals(claimsKey, claims)
		return c.Next()
	}
}

// RequireRole restricts the route to the given roles (must run after RequireAuth).
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals(claimsKey).(*auth.Claims)
		if !ok || claims == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthenticated")
		}
		for _, r := range roles {
			if claims.Role == r {
				return c.Next()
			}
		}
		return fiber.NewError(fiber.StatusForbidden, "insufficient permissions")
	}
}

func claimsFrom(c *fiber.Ctx) *auth.Claims {
	if v, ok := c.Locals(claimsKey).(*auth.Claims); ok {
		return v
	}
	return nil
}

// GetOrgID returns the tenant org id derived from the authenticated principal.
func GetOrgID(c *fiber.Ctx) uuid.UUID {
	if cl := claimsFrom(c); cl != nil {
		return cl.OrgID
	}
	return uuid.Nil
}

func GetAccountID(c *fiber.Ctx) uuid.UUID {
	if cl := claimsFrom(c); cl != nil {
		return uuid.MustParse(cl.Subject)
	}
	return uuid.Nil
}

func GetPersonID(c *fiber.Ctx) uuid.UUID {
	if cl := claimsFrom(c); cl != nil {
		return cl.PersonID
	}
	return uuid.Nil
}

func GetRole(c *fiber.Ctx) string {
	if cl := claimsFrom(c); cl != nil {
		return cl.Role
	}
	return ""
}
