//go:build wireinject
// +build wireinject

package main

import (
	"fmt"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/config"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/database"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/handlers"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/middlewares"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/repository"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/routes"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Application struct {
	App    *fiber.App
	Logger *logrus.Logger
}

// =====================================================================================
// SECTION: Infrastructure Providers
// =====================================================================================

func ProvideLogger() *logrus.Logger {
	return config.InitLogger()
}

func ProvideValidator() *validator.Validate {
	return config.InitValidator()
}

func ProvideDatabase(logger *logrus.Logger) (*gorm.DB, error) {
	db, err := database.InitDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	logger.Info("Database connection and migration successful")
	return db, nil
}

func ProvideRedis(logger *logrus.Logger) (*redis.Client, error) {
	client, err := database.InitRedis()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}
	logger.Info("Redis connection successful")
	return client, nil
}

func ProvideRouter(h *handlers.Handlers, middleware *middlewares.Middleware) *fiber.App {
	app := fiber.New()
	routes.RegisterRoutes(app, h, middleware)
	return app
}

func NewApplication(app *fiber.App, logger *logrus.Logger) Application {
	return Application{
		App:    app,
		Logger: logger,
	}
}

// =====================================================================================
// SECTION: Wire Sets - Organized by Layer
// =====================================================================================

var InfrastructureSet = wire.NewSet(
	ProvideLogger,
	ProvideValidator,
	ProvideDatabase,
	ProvideRedis,
)

var RepositorySet = wire.NewSet(
	repository.NewAuthRepository,
	repository.NewSessionRepository,
	repository.NewMahasiswaRepository,
	repository.NewTahunAkademikRepository,
	repository.NewDosenRepository,
	repository.NewPegawaiRepository,
)

var ServiceSet = wire.NewSet(
	service.NewAuthService,
	service.NewMahasiswaService,
	service.NewTahunAkademikService,
)

var HandlerSet = wire.NewSet(
	handlers.NewAuthHandler,
	handlers.NewMahasiswaHandler,
	handlers.NewHandlers,
)

var MiddlewareSet = wire.NewSet(
	middlewares.NewMiddleware,
)

var RouterSet = wire.NewSet(
	ProvideRouter,
)

// =====================================================================================
// SECTION: Main Application Set
// =====================================================================================

var AppSet = wire.NewSet(
	InfrastructureSet,
	RepositorySet,
	ServiceSet,
	HandlerSet,
	MiddlewareSet,
	RouterSet,
	NewApplication,
)

func InitializeApp() (Application, error) {
	wire.Build(AppSet)
	return Application{}, nil
}
