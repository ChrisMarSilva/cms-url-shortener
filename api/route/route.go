package route

import (
	"time"

	"github.com/ChrisMarSilva/cms-url-shortener/handlers"
	"github.com/ChrisMarSilva/cms-url-shortener/repositories"
	"github.com/ChrisMarSilva/cms-url-shortener/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"

	// "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/timeout"
)

func NewRoutes() *fiber.App {

	defaultHandler := handlers.NewDefaultHandler()

	resolveRepo := repositories.NewResolveRepository()
	resolveServ := services.NewResolveService(*resolveRepo)
	resolveHandler := handlers.NewResolveHandler(*resolveServ)

	shortenRepo := repositories.NewShortenRepository()
	shortenServ := services.NewShortenService(*shortenRepo)
	shortenHandler := handlers.NewShortenHandler(*shortenServ)

	app := fiber.New()

	app.Use(compress.New(compress.Config{Level: compress.LevelBestSpeed}))

	//app.Use(cors.New())
	app.Use(cors.New(cors.Config{AllowOrigins: "*", AllowMethods: "*", AllowHeaders: "*", AllowCredentials: true}))

	app.Use(encryptcookie.New(encryptcookie.Config{Key: "7tDyMRLm2ii3BVDiN7GXfKoALsiMMzrr"}))

	// app.Use(logger.New())
	// app.Use(logger.New(logger.Config{Format: "[${time}]: ${ip} ${status} ${latency} ${method} ${path} Error: ${error}\n"}))

	app.Use(recover.New())

	// app.Get("/:url", resolveHandler.ResolveURL)
	app.Get("/:url", timeout.New(resolveHandler.ResolveURL, 10*time.Second))

	app.Post("/api/v1", shortenHandler.ShortenURL)
	// app.Post("/api/v1", timeout.New(shortenHandler.ShortenURL, 10*time.Second))

	// app.Get("/api/v1/dashboard", monitor.New())

	app.Use(defaultHandler.NotFound)

	return app
}
