package handlers

import (
	"log"

	"github.com/ChrisMarSilva/cms-url-shortener/databases"
	"github.com/ChrisMarSilva/cms-url-shortener/services"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
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
	log.Println(url)
	// return handler.service.ResolveURL(url)

	r := databases.CreateClient(0)
	defer r.Close()

	value, err := r.Get(databases.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short not found in the database"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot connext to DB"})
	}

	rInr := databases.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(databases.Ctx, "counter")

	log.Println(value)
	return c.Redirect(value, fiber.StatusMovedPermanently)
}
