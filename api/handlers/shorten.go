package handlers

import (
	"os"

	// "strconv"
	// "time"

	//"github.com/ChrisMarSilva/cms-url-shortener/databases"
	"github.com/ChrisMarSilva/cms-url-shortener/entities"
	"github.com/ChrisMarSilva/cms-url-shortener/helpers"
	"github.com/ChrisMarSilva/cms-url-shortener/services"
	"github.com/asaskevich/govalidator"

	// "github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
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

	// ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	// defer cancel()
	// select {
	// case <-ctx.Done():
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Sem Context"})
	// 	break
	// }

	body := new(entities.Request)
	// return handler.service.ShortenURL(body)

	if err := c.BodyParser(&body); err != nil {
		// log.Println("BodyParser", string(c.Body()), err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"}) // err.Error()
	}

	// implement rate limiting

	// r2 := databases.CreateClient(1)
	// defer r2.Close()

	// my_ip := c.IP()

	// val, err := r2.Get(databases.Ctx, my_ip).Result()
	// if err == redis.Nil {
	// 	val = os.Getenv("API_QUOTA")
	// 	_ = r2.Set(databases.Ctx, my_ip, val, 30*60*time.Second).Err()
	// 	// } else {
	// 	// 	val, _ = r2.Get(databases.Ctx, my_ip).Result()
	// 	// 	valInt, _ := strconv.Atoi(val)
	// 	// 	if valInt <= 0 {
	// 	// 		limit, _ := r2.TTL(databases.Ctx, my_ip).Result()
	// 	// 		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Rate limit exceded", "rate_limit_rest": limit / time.Nanosecond / time.Minute})
	// 	// 	}
	// }
	// valInt, _ := strconv.Atoi(val)
	// if valInt <= 0 {
	// 	limit, _ := r2.TTL(databases.Ctx, my_ip).Result()
	// 	return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Rate limit exceded", "rate_limit_rest": limit / time.Nanosecond / time.Minute})
	// }

	// check if the input if an actual URL

	if !govalidator.IsURL(body.URL) {
		// log.Println("IsURL.body.URL", body.URL, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	// check for domain error

	if !helpers.RemoveDomainError(body.URL) {
		// log.Println("RemoveDomainError.body.URL", body.URL, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "you can't hack the system (:"})
	}

	// enforce https, SSL

	body.URL = helpers.EnforceHTTP(body.URL)

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	id := body.CustomShort
	if id == "" {
		id = uuid.New().String()[:6]
	}

	body.CustomShort = id

	// r := databases.CreateClient(0)
	// defer r.Close()

	// val, _ = r.Get(databases.Ctx, id).Result()
	// if val == "" {
	// 	// return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "URL custom short is already in use"})
	// 	err = r.Set(databases.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to connext to server"})
	// 	}
	// }

	resp := entities.Response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	// r2.Decr(databases.Ctx, my_ip)

	// val, _ = r2.Get(databases.Ctx, my_ip).Result()
	// resp.XRateRemaining, _ = strconv.Atoi(val)

	// ttl, _ := r2.TTL(databases.Ctx, my_ip).Result()
	// resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("API_DOMAIN") + "/" + id

	conn, err := amqp.Dial(os.Getenv("RABBIT_MQ_ADDR"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to connect to RabbitMQ"})
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open a channel"})
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(os.Getenv("RABBIT_MQ_QUEUE"), true, false, false, false, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to declare a queue"})
	}

	// time.Sleep(10 * time.Second)
	// ctx := c.Context()
	// log.Println(ctx)

	if c.Context() == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Sem Context"})
	}

	// enviar a url, ip, id
	message := amqp.Publishing{ContentType: "text/plain", Body: []byte(body.URL)}

	// err = ch.Publish("", q.Name, false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte(body.URL)})
	err = ch.Publish("", q.Name, false, false, message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to publish a message"})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
