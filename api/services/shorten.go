package services

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ChrisMarSilva/cms-url-shortener/databases"
	"github.com/ChrisMarSilva/cms-url-shortener/entities"
	"github.com/ChrisMarSilva/cms-url-shortener/helpers"
	"github.com/ChrisMarSilva/cms-url-shortener/repositories"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	opentracing "github.com/opentracing/opentracing-go"
	spanlog "github.com/opentracing/opentracing-go/log"
	"github.com/streadway/amqp"
)

type ShortenService struct {
	repo     repositories.ShortenRepository
	rabbitmq helpers.IRabbitMQ
}

func NewShortenService(repo repositories.ShortenRepository, rabbitmq helpers.IRabbitMQ) *ShortenService {
	return &ShortenService{
		repo:     repo,
		rabbitmq: rabbitmq,
	}
}

func (service *ShortenService) ShortenURL(ctx context.Context, body *entities.Request, my_ip string) *entities.HttpResponse {

	sp, _ := opentracing.StartSpanFromContext(ctx, "ShortenService.ShortenURL")
	defer sp.Finish()

	sp.SetTag("body", body)

	// check if the input if an actual URL
	if !govalidator.IsURL(body.URL) {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(errors.New("Invalid URL")))
		return entities.BadRequest("Invalid URL")
	}

	// check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(errors.New("you can't hack the system (:")))
		return entities.BadRequest("you can't hack the system (:")
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

	resp := entities.ShortenResponse{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	resp.CustomShort = os.Getenv("API_DOMAIN") + "/" + id
	sp.SetTag("resp", resp)

	sp3 := opentracing.StartSpan("amqp.Dial", opentracing.ChildOf(sp.Context()))
	conn, err := amqp.Dial(os.Getenv("RABBIT_MQ_ADDR"))
	sp3.Finish()
	if err != nil {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(err))
		return entities.InternalServerError("Failed to connect to RabbitMQ")
	}
	defer conn.Close()

	sp4 := opentracing.StartSpan("amqp.Channel", opentracing.ChildOf(sp.Context()))
	ch, err := conn.Channel()
	sp4.Finish()
	if err != nil {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(err))
		return entities.InternalServerError("Failed to open a channel")
	}
	defer ch.Close()

	sp5 := opentracing.StartSpan("amqp.QueueDeclare", opentracing.ChildOf(sp.Context()))
	q, err := ch.QueueDeclare(os.Getenv("RABBIT_MQ_QUEUE"), true, false, false, false, nil)
	sp5.Finish()
	if err != nil {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(err))
		return entities.InternalServerError("Failed to declare a queue")
	}

	if ctx == nil {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(errors.New("Sem Context")))
		return entities.InternalServerError("Sem Context")
	}

	// enviar a url, ip, id
	// "application/json"
	message := amqp.Publishing{DeliveryMode: amqp.Persistent, ContentType: "text/plain", Body: []byte(body.URL)}

	sp6 := opentracing.StartSpan("amqp.Publish", opentracing.ChildOf(sp.Context()))
	// err = ch.Publish("", q.Name, false, false, amqp.Publishing{DeliveryMode: amqp.Persistent, ContentType: "text/plain", Body: []byte(body.URL)})
	err = ch.Publish("", q.Name, false, false, message)
	sp6.Finish()
	if err != nil {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(err))
		return entities.InternalServerError("Failed to publish a message")
	}

	return entities.Ok(resp)
}

func (service *ShortenService) ShortenURL2(ctx context.Context, body *entities.Request, my_ip string) *entities.HttpResponse {

	sp, _ := opentracing.StartSpanFromContext(ctx, "ShortenService.ShortenURL2")
	defer sp.Finish()

	sp.SetTag("body", body)

	// check if the input if an actual URL
	if !govalidator.IsURL(body.URL) {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(errors.New("Invalid URL")))
		return entities.BadRequest("Invalid URL")
	}

	// check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(errors.New("you can't hack the system (:")))
		return entities.BadRequest("you can't hack the system (:")
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

	resp := entities.ShortenResponse{
		URL:             body.URL,
		CustomShort:     os.Getenv("API_DOMAIN") + "/" + id,
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}
	sp.SetTag("resp", resp)

	sp3 := opentracing.StartSpan("rabbitmq.SendMessage", opentracing.ChildOf(sp.Context()))
	if err := service.rabbitmq.SendMessage(os.Getenv("RABBIT_MQ_QUEUE"), body.URL); err != nil {
		sp3.Finish()
		log.Println(err)
		//  Reason: "channel/connection is not open"
		//  Exception (504) Reason: "channel/connection is not open"
		// "UNEXPECTED_FRAME - expected content header for class 60, got non content header frame instead"
		return entities.InternalServerError("Failed to connect to RabbitMQ")
	}
	sp3.Finish()

	return entities.Ok(resp)
}

func (service *ShortenService) ShortenURL3(ctx context.Context, body *entities.Request, my_ip string) *entities.HttpResponse {

	sp, _ := opentracing.StartSpanFromContext(ctx, "ShortenService.ShortenURL")
	defer sp.Finish()

	sp.SetTag("body", body)

	// implement rate limiting

	sp2 := opentracing.StartSpan("r2.CreateClient", opentracing.ChildOf(sp.Context()))
	r2 := databases.CreateClient(1)
	sp2.Finish()
	defer r2.Close()

	sp3 := opentracing.StartSpan("r2.Get.my_ip", opentracing.ChildOf(sp.Context()))
	val, err := r2.Get(databases.Ctx, my_ip).Result()
	sp3.Finish()
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
		// limit, _ := r2.TTL(databases.Ctx, my_ip).Result()
		return entities.InternalServerError("Rate limit exceded. rate_limit_rest: ") // limit / time.Nanosecond / time.Minute
	}

	// check if the input if an actual URL
	if !govalidator.IsURL(body.URL) {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(errors.New("Invalid URL")))
		return entities.BadRequest("Invalid URL")
	}

	// check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(errors.New("you can't hack the system (:")))
		return entities.BadRequest("you can't hack the system (:")
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
		sp4 := opentracing.StartSpan("UUID", opentracing.ChildOf(sp.Context()))
		id = uuid.New().String()[:6]
		sp4.Finish()
	}
	sp.SetTag("body.id", id)

	sp5 := opentracing.StartSpan("r.databases.CreateClient", opentracing.ChildOf(sp.Context()))
	r := databases.CreateClient(0)
	sp5.Finish()
	defer r.Close()

	sp6 := opentracing.StartSpan("r.Get.id", opentracing.ChildOf(sp.Context()))
	val, _ = r.Get(databases.Ctx, id).Result()
	sp6.Finish()
	if val == "" {
		// return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "URL custom short is already in use"})
		err = r.Set(databases.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
		if err != nil {
			return entities.InternalServerError("Unable to connext to server")
		}
	}

	resp := entities.ShortenResponse{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}
	sp.SetTag("resp", resp)

	sp7 := opentracing.StartSpan("r2.Decr.my_ip", opentracing.ChildOf(sp.Context()))
	r2.Decr(databases.Ctx, my_ip)
	sp7.Finish()

	sp8 := opentracing.StartSpan("r2.Get.my_ip", opentracing.ChildOf(sp.Context()))
	val, _ = r2.Get(databases.Ctx, my_ip).Result()
	sp8.Finish()

	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(databases.Ctx, my_ip).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("API_DOMAIN") + "/" + id

	return entities.Ok(resp)
}
