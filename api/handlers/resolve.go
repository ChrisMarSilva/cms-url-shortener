package handlers

import (
	"github.com/ChrisMarSilva/cms-url-shortener/services"
	"github.com/gofiber/fiber/v2"
	opentracing "github.com/opentracing/opentracing-go"
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

	url := c.Params("url")

	sp, ctx := opentracing.StartSpanFromContext(c.Context(), "ResolveHandler.ResolveURL")
	sp.SetTag("Method", c.Method())
	sp.SetTag("URL", c.Path())
	sp.SetTag("Params.url", url)
	defer sp.Finish()

	response := handler.service.ResolveURL(ctx, url)

	if response.StatusCode == fiber.StatusOK {
		return c.Redirect(response.URL, fiber.StatusMovedPermanently)
	}

	return c.Status(response.StatusCode).JSON(response.Data)
}
