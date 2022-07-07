package handlers

import (
	"github.com/ChrisMarSilva/cms-url-shortener/services"
	"github.com/gofiber/fiber/v2"
	opentracing "github.com/opentracing/opentracing-go"
)

type DefaultHandler struct {
	service services.DefaultService
}

func NewDefaultHandler(service services.DefaultService) *DefaultHandler {
	return &DefaultHandler{
		service: service,
	}
}

func (handler *DefaultHandler) NotFound(c *fiber.Ctx) error {

	sp, ctx := opentracing.StartSpanFromContext(c.Context(), "DefaultHandler.NotFound")
	sp.SetTag("Method", c.Method())
	sp.SetTag("URL", c.Path())
	defer sp.Finish()

	response := handler.service.NotFound(ctx)
	return c.Status(response.StatusCode).JSON(response.Data)
}
