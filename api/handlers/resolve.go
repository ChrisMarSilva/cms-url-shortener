package handlers

import (
	"github.com/ChrisMarSilva/cms-url-shortener/databases"
	"github.com/ChrisMarSilva/cms-url-shortener/services"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	opentracing "github.com/opentracing/opentracing-go"
	spanlog "github.com/opentracing/opentracing-go/log"
)

type ResolveHandler struct {
	service services.ResolveService
}

func NewResolveHandler(service services.ResolveService) *ResolveHandler {
	return &ResolveHandler{
		service: service,
	}
}

func (handler *ResolveHandler) ResolveURL(c *fiber.Ctx) error {

	sp, _ := opentracing.StartSpanFromContext(c.Context(), "ResolveHandler.ResolveURL")
	sp.SetTag("Method", c.Method())
	sp.SetTag("URL", c.Path())
	defer sp.Finish()

	url := c.Params("url")
	sp.SetTag("Params.url", url)

	// return handler.service.ResolveURL(url)

	sp2 := opentracing.StartSpan("BD.Conn", opentracing.ChildOf(sp.Context()))
	r := databases.CreateClient(0)
	sp2.Finish()
	defer r.Close()

	sp3 := opentracing.StartSpan("BD.Get", opentracing.ChildOf(sp.Context()))
	value, err := r.Get(databases.Ctx, url).Result()
	sp3.Finish()
	if err == redis.Nil {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short not found in the database"})
	} else if err != nil {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot connext to DB"})
	}
	sp.SetTag("value", value)

	// rInr := databases.CreateClient(1)
	// defer rInr.Close()

	// _ = rInr.Incr(databases.Ctx, "counter")

	return c.Redirect(value, fiber.StatusMovedPermanently)
}
