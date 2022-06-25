package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type DefaultHandler struct {
}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{}
}

func (handler *DefaultHandler) NotFound(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotFound)
}
