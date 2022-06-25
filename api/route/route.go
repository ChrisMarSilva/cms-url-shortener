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
	"github.com/gofiber/fiber/v2/middleware/logger"
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

	app.Use(encryptcookie.New(encryptcookie.Config{Key: "VWY5XBVap84Zpd0ckbT1reTl0NM6pz7R"}))

	//app.Use(logger.New())
	app.Use(logger.New(logger.Config{Format: "[${time}]: ${ip} ${status} ${latency} ${method} ${path} Error: ${error}\n"}))

	app.Use(recover.New())

	// app.Get("/:url", routes.ResolveURL)
	app.Get("/:url", timeout.New(resolveHandler.ResolveURL, 5*time.Second))

	// app.Post("/api/v1", routes.ShortenURL)
	app.Post("/api/v1", timeout.New(shortenHandler.ShortenURL, 5*time.Second))

	// app.Get("/api/v1/dashboard", monitor.New())
	app.Use(defaultHandler.NotFound)

	return app
}
