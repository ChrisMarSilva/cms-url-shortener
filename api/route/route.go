package route

import (
	"time"

	"github.com/ChrisMarSilva/cms-url-shortener/handlers"
	"github.com/ChrisMarSilva/cms-url-shortener/repositories"
	"github.com/ChrisMarSilva/cms-url-shortener/services"
	"github.com/aschenmaker/fiber-opentracing/fjaeger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	jconfig "github.com/uber/jaeger-client-go/config"
)

func NewRoutes() *fiber.App {

	app := fiber.New()

	app.Use(compress.New(compress.Config{Level: compress.LevelBestSpeed}))
	app.Use(cors.New(cors.Config{AllowOrigins: "*", AllowMethods: "*", AllowHeaders: "*", AllowCredentials: true}))
	app.Use(encryptcookie.New(encryptcookie.Config{Key: "7tDyMRLm2ii3BVDiN7GXfKoALsiMMzrr"}))
	app.Use(logger.New(logger.Config{Format: "[${time}]: ${ip} ${status} ${latency} ${method} ${path} Error: ${error}\n"}))
	app.Use(recover.New())

	cfg := fjaeger.Config{
		ServiceName:      "cms.url.shortener",
		Sampler:          &jconfig.SamplerConfig{Type: "const", Param: 1},
		Reporter:         &jconfig.ReporterConfig{LogSpans: true, BufferFlushInterval: 1 * time.Second},
		EnableRPCMetrics: true,
		PanicOnError:     false,
	}
	fjaeger.New(cfg)

	// log.Println("cfg", cfg)
	// log.Println("Sampler", cfg.Sampler)
	// log.Println("SamplingServerURL", cfg.Sampler.SamplingServerURL)
	// log.Println("Getenv", os.Getenv("JAEGER_SAMPLING_ENDPOINT"))

	resolveRepo := repositories.NewResolveRepository()
	resolveServ := services.NewResolveService(*resolveRepo)
	resolveHandler := handlers.NewResolveHandler(*resolveServ)
	app.Get("/:url", timeout.New(resolveHandler.ResolveURL, 10*time.Second)) // resolveHandler.ResolveURL // timeout.New(resolveHandler.ResolveURL, 10*time.Second)

	shortenRepo := repositories.NewShortenRepository()
	shortenServ := services.NewShortenService(*shortenRepo)
	shortenHandler := handlers.NewShortenHandler(*shortenServ)
	app.Post("/api/v1", shortenHandler.ShortenURL) // shortenHandler.ShortenURL // timeout.New(shortenHandler.ShortenURL, 10*time.Second)

	defaultHandler := handlers.NewDefaultHandler()
	app.Use(defaultHandler.NotFound)

	return app
}
