package handlers

import (
	"github.com/ChrisMarSilva/cms-url-shortener/entities"
	"github.com/ChrisMarSilva/cms-url-shortener/services"

	"github.com/gofiber/fiber/v2"
	opentracing "github.com/opentracing/opentracing-go"
	spanlog "github.com/opentracing/opentracing-go/log"
)

type ShortenHandler struct {
	service services.ShortenService
}

func NewShortenHandler(service services.ShortenService) *ShortenHandler {
	return &ShortenHandler{
		service: service,
	}
}

func (handler *ShortenHandler) ShortenURL(c *fiber.Ctx) error {

	sp, ctx := opentracing.StartSpanFromContext(c.Context(), "ShortenHandler.ShortenURL")
	sp.SetTag("Method", c.Method())
	sp.SetTag("URL", c.Path())
	defer sp.Finish()

	body := new(entities.Request)

	//sp2 := opentracing.StartSpan("BodyParser", opentracing.ChildOf(sp.Context()))
	if err := c.BodyParser(&body); err != nil {
		//sp2.Finish()
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}
	//sp2.Finish()

	// time.Sleep(10 * time.Second)
	// ctx := c.Context()
	// log.Println(ctx)

	my_ip := c.IP()

	response := handler.service.ShortenURL3(ctx, body, my_ip) // ShortenURL // ShortenURL2 // ShortenURL3

	return c.Status(response.StatusCode).JSON(response.Data)
}
