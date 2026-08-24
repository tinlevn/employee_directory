package main

import (
	"context"
	"log"
	"os"
	"time"

	"employee-directory-api/internal/config"
	"employee-directory-api/internal/handler"
	"employee-directory-api/internal/middleware"
	"employee-directory-api/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()
	validate := validator.New(validator.WithRequiredStructEnabled())

	// DB pool is optional for dev - if DATABASE_URL not reachable, run in no-db mode for health/docs
	var pool *pgxpool.Pool
	if cfg.DatabaseURL != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		p, err := repository.NewPool(ctx, cfg.DatabaseURL)
		if err != nil {
			log.Printf("WARN: db connect failed (%v) - running without DB (health will be degraded)", err)
		} else {
			pool = p
			defer pool.Close()
			log.Println("connected to postgres")
		}
	}

	app := fiber.New(fiber.Config{
		ErrorHandler:          middleware.ErrorHandler,
		BodyLimit:             1 * 1024 * 1024,
		DisableStartupMessage: false,
		AppName:               "employee-directory-api",
	})

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{Format: "${time} | ${status} | ${latency} | ${ip} | ${method} ${path} | ${locals:requestid}\n"}))
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     joinOrigins(cfg.AllowedOrigins),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Request-Id",
		AllowMethods:     "GET,POST,PATCH,DELETE,OPTIONS",
		AllowCredentials: false,
	}))

	// health
	app.Get("/health", func(c *fiber.Ctx) error {
		status := "healthy"
		code := fiber.StatusOK
		dbStatus := "connected"
		if pool == nil {
			dbStatus = "disconnected"
			status = "degraded"
		} else {
			ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
			defer cancel()
			if err := pool.Ping(ctx); err != nil {
				dbStatus = "unreachable"
				status = "unhealthy"
				code = fiber.StatusServiceUnavailable
			}
		}
		return c.Status(code).JSON(fiber.Map{
			"status":    status,
			"database":  dbStatus,
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})

	api := app.Group("/api/v1")

	if pool != nil {
		// repos
		personRepo := repository.NewPersonRepository(pool)
		emergencyRepo := repository.NewEmergencyContactRepository(pool)
		employmentRepo := repository.NewEmploymentRepository(pool)
		eventRepo := repository.NewEventRepository(pool)
		transferRepo := repository.NewTransferRepository(pool)
		analyticsRepo := repository.NewAnalyticsRepository(pool)
		orgRepo := repository.NewOrgRepository(pool)

		// handlers
		personH := handler.NewPersonHandler(personRepo, validate)
		emergencyH := handler.NewEmergencyHandler(emergencyRepo, validate)
		employmentH := handler.NewEmploymentHandler(employmentRepo, validate)
		eventH := handler.NewEventHandler(eventRepo, transferRepo, validate)
		analyticsH := handler.NewAnalyticsHandler(analyticsRepo, validate)
		orgH := handler.NewOrgHandler(orgRepo, validate)

		// orgs
		api.Get("/organizations", orgH.List)
		api.Post("/organizations", orgH.Create)
		api.Get("/organizations/:id", orgH.Get)

		// persons
		api.Get("/persons", personH.List)
		api.Post("/persons", personH.Create)
		api.Get("/persons/:id", personH.Get)
		api.Patch("/persons/:id", personH.Update)
		api.Delete("/persons/:id", personH.Delete)

		// emergency contact — separate entity, optional FK
		api.Get("/persons/:id/emergency-contact", emergencyH.Get)
		api.Post("/persons/:id/emergency-contact", emergencyH.Upsert)
		api.Patch("/persons/:id/emergency-contact", emergencyH.Update)
		api.Delete("/persons/:id/emergency-contact", emergencyH.Delete)

		// employment
		api.Get("/persons/:id/employment", employmentH.List)
		api.Get("/persons/:id/employment/current", employmentH.GetCurrent)
		api.Post("/persons/:id/employment", employmentH.Create)

		// events
		api.Get("/persons/:id/events", eventH.ListEvents)
		api.Post("/persons/:id/events", eventH.CreateEvent)

		// transfers
		api.Get("/persons/:id/transfers", eventH.ListTransfers)
		api.Post("/persons/:id/transfers", eventH.CreateTransfer)

		// analytics
		api.Get("/analytics/headcount", analyticsH.Headcount)
		api.Get("/analytics/attrition", analyticsH.Attrition)
		api.Get("/analytics/movements", analyticsH.Movements)
		api.Get("/analytics/snapshot/:date", analyticsH.Snapshot)

		// compat: keep old /api/employees -> redirect to /persons for old Angular client
		app.Get("/api/employees", func(c *fiber.Ctx) error {
			c.Path("/api/v1/persons")
			return personH.List(c)
		})
		app.Get("/api/employees/:id", func(c *fiber.Ctx) error {
			c.Params("id", c.Params("id"))
			return personH.Get(c)
		})
	} else {
		api.All("/*", func(c *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusServiceUnavailable, "database not connected - set DATABASE_URL")
		})
	}

	// OpenAPI document is kept as a checked-in artifact so clients can inspect the actual contract.
	app.Get("/openapi.yaml", func(c *fiber.Ctx) error {
		for _, path := range []string{"openapi.yaml", "/openapi.yaml", "../openapi.yaml"} {
			if document, err := os.ReadFile(path); err == nil {
				c.Set(fiber.HeaderContentType, "application/yaml; charset=utf-8")
				return c.Send(document)
			}
		}
		return fiber.NewError(fiber.StatusNotFound, "OpenAPI document unavailable")
	})
	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Redirect("/openapi.yaml", fiber.StatusMovedPermanently)
	})

	app.Use(middleware.NotFoundHandler)

	log.Printf("listening on :%s (env=%s)", cfg.Port, cfg.Env)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

func joinOrigins(origins []string) string {
	if len(origins) == 0 {
		return "*"
	}
	out := ""
	for i, o := range origins {
		if i > 0 {
			out += ", "
		}
		out += o
	}
	return out
}
