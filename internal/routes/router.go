package routes

import (
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/handlers"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/middlewares"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func RegisterRoutes(app *fiber.App, h *handlers.Handlers, mid *middlewares.Middleware) {
	app.Use(logger.New())

	app.Use(
		helmet.New(),
	)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:4173,https://krs-dev.masako.my.id",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH",
		AllowCredentials: true,
	}))

	middlewares.InitTimezoneWib()

	app.Use(limiter.New(limiter.Config{
		Max:               50,
		Expiration:        1 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  fiber.StatusTooManyRequests,
				"message": "Terlalu banyak permintaan, coba lagi nanti.",
			})
		},
	}))

	api := app.Group("/api/v1")
	api.Use(mid.KRSScheduleMiddleware())

	authRoute := api.Group("/auth")
	authRoute.Post("/login", h.AuthHandler.Login)
	authRoute.Post("/logout", mid.AuthMiddleware(), h.AuthHandler.Logout)
	authRoute.Get("/session", mid.AuthMiddleware(), h.AuthHandler.GetSession)

	mahasiswaRoute := api.Group("/mahasiswa", mid.AuthMiddleware())
	mahasiswaRoute.Get("/syarat-pengisian-krs", h.MahasiswaHandler.GetSyaratPengisisanKRS)
	mahasiswaRoute.Get("/informasi-umum", h.MahasiswaHandler.GetInformasiUmum)
	mahasiswaRoute.Get("/penawaran-kelas", h.MahasiswaHandler.GetPenawaranKelas)
	mahasiswaRoute.Post("/status-kouta-kelas", h.MahasiswaHandler.GetStatusKoutaKelas)
	mahasiswaRoute.Post("/status-kouta-kelas-batch", h.MahasiswaHandler.GetStatusKoutaKelasBatch)
}
