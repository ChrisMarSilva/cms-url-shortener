package handlers

import (
	"errors"
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
	opentracing "github.com/opentracing/opentracing-go"
	spanlog "github.com/opentracing/opentracing-go/log"
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

	sp, _ := opentracing.StartSpanFromContext(c.Context(), "ShortenHandler.ShortenURL")
	sp.SetTag("Method", c.Method())
	sp.SetTag("URL", c.Path())
	defer sp.Finish()

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
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"}) // err.Error()
	}

	sp.SetTag("body", body)
	sp.SetTag("body.URL", body.URL)
	sp.SetTag("body.CustomShort", body.CustomShort)
	sp.SetTag("body.Expiry", body.Expiry)

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
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(errors.New("Invalid URL")))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	// check for domain error

	if !helpers.RemoveDomainError(body.URL) {
		// log.Println("RemoveDomainError.body.URL", body.URL, err)
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(errors.New("you can't hack the system (:")))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "you can't hack the system (:"})
	}

	// enforce https, SSL

	body.URL = helpers.EnforceHTTP(body.URL)
	sp.SetTag("body.URL", body.URL)

	if body.Expiry == 0 {
		body.Expiry = 24
	}
	sp.SetTag("body.Expiry", body.Expiry)

	id := body.CustomShort
	if id == "" {
		sp2 := opentracing.StartSpan("UUID", opentracing.ChildOf(sp.Context()))
		id = uuid.New().String()[:6]
		sp2.Finish()
	}
	sp.SetTag("body.id", id)

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

	sp3 := opentracing.StartSpan("amqp.Dial", opentracing.ChildOf(sp.Context()))
	conn, err := amqp.Dial(os.Getenv("RABBIT_MQ_ADDR"))
	sp3.Finish()
	if err != nil {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to connect to RabbitMQ"})
	}
	defer conn.Close()

	sp4 := opentracing.StartSpan("amqp.Channel", opentracing.ChildOf(sp.Context()))
	ch, err := conn.Channel()
	sp4.Finish()
	if err != nil {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open a channel"})
	}
	defer ch.Close()

	sp5 := opentracing.StartSpan("amqp.QueueDeclare", opentracing.ChildOf(sp.Context()))
	q, err := ch.QueueDeclare(os.Getenv("RABBIT_MQ_QUEUE"), true, false, false, false, nil)
	sp5.Finish()
	if err != nil {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to declare a queue"})
	}

	// time.Sleep(10 * time.Second)
	// ctx := c.Context()
	// log.Println(ctx)

	if c.Context() == nil {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(errors.New("Sem Context")))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Sem Context"})
	}

	// enviar a url, ip, id
	message := amqp.Publishing{ContentType: "text/plain", Body: []byte(body.URL)}

	sp6 := opentracing.StartSpan("amqp.Publish", opentracing.ChildOf(sp.Context()))
	// err = ch.Publish("", q.Name, false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte(body.URL)})
	err = ch.Publish("", q.Name, false, false, message)
	sp6.Finish()
	if err != nil {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to publish a message"})
	}

	sp.SetTag("resp", resp)
	sp.SetTag("resp.URL", resp.URL)
	sp.SetTag("resp.CustomShort", resp.CustomShort)
	sp.SetTag("resp.Expiry", resp.Expiry)
	sp.SetTag("resp.XRateRemaining", resp.XRateRemaining)
	sp.SetTag("resp.XRateLimitReset", resp.XRateLimitReset)

	return c.Status(fiber.StatusOK).JSON(resp)
}
