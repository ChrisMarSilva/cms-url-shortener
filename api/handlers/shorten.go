package handlers

import (
	"os"
	"strconv"
	"time"

	"github.com/ChrisMarSilva/cms-url-shortener/databases"
	"github.com/ChrisMarSilva/cms-url-shortener/entities"
	"github.com/ChrisMarSilva/cms-url-shortener/helpers"
	"github.com/ChrisMarSilva/cms-url-shortener/services"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	body := new(entities.Request)
	// return handler.service.ShortenURL(body)

	if err := c.BodyParser(&body); err != nil {
		// log.Println(string(c.Body()))
		// log.Println("c.BodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"}) // err.Error()
	}

	// implement rate limiting

	r2 := databases.CreateClient(1)
	defer r2.Close()

	my_ip := c.IP()

	val, err := r2.Get(databases.Ctx, my_ip).Result()
	if err == redis.Nil {
		val = os.Getenv("API_QUOTA")
		_ = r2.Set(databases.Ctx, my_ip, val, 30*60*time.Second).Err()
		// } else {
		// 	val, _ = r2.Get(databases.Ctx, my_ip).Result()
		// 	valInt, _ := strconv.Atoi(val)
		// 	if valInt <= 0 {
		// 		limit, _ := r2.TTL(databases.Ctx, my_ip).Result()
		// 		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Rate limit exceded", "rate_limit_rest": limit / time.Nanosecond / time.Minute})
		// 	}
	}
	valInt, _ := strconv.Atoi(val)
	if valInt <= 0 {
		limit, _ := r2.TTL(databases.Ctx, my_ip).Result()
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Rate limit exceded", "rate_limit_rest": limit / time.Nanosecond / time.Minute})
	}

	// check if the input if an actual URL

	if !govalidator.IsURL(body.URL) {
		// log.Println("IsURL.body.URL", body.URL)
		// log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	// check for domain error

	if !helpers.RemoveDomainError(body.URL) {
		// log.Println("RemoveDomainError.body.URL", body.URL)
		// log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "you can't hack the system (:"})
	}

	// enforce https, SSL

	body.URL = helpers.EnforceHTTP(body.URL)

	// var id string
	id := body.CustomShort

	if id == "" {
		id = uuid.New().String()[:6]
		// } else {
		//	id = body.CustomShort
	}

	r := databases.CreateClient(0)
	defer r.Close()

	val, _ = r.Get(databases.Ctx, id).Result()
	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "URL custom short is already in use"})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = r.Set(databases.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to connext to server"})
	}

	resp := entities.Response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	r2.Decr(databases.Ctx, my_ip)

	val, _ = r2.Get(databases.Ctx, my_ip).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(databases.Ctx, my_ip).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("API_DOMAIN") + "/" + id

	return c.Status(fiber.StatusOK).JSON(resp)
}
