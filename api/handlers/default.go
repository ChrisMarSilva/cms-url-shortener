package handlers

import (
	"github.com/gofiber/fiber/v2"
	opentracing "github.com/opentracing/opentracing-go"
)

type DefaultHandler struct {
}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{}
}

func (handler *DefaultHandler) NotFound(c *fiber.Ctx) error {

	sp, _ := opentracing.StartSpanFromContext(c.Context(), "DefaultHandler.NotFound")
	sp.SetTag("Method", c.Method())
	sp.SetTag("URL", c.Path())
	defer sp.Finish()

	// time.Sleep(1 * time.Second)

	// sp.LogFields(spanlog.String("event", "soft error"), spanlog.String("type", "cache timeout"), spanlog.Int("waited.millis", 1500))
	// sp.LogFields(spanlog.String("event", "getResponse"), spanlog.String("value", "string(body)"))
	// sp.LogEvent("AAA")
	// sp.SetBaggageItem("error3", "123")
	// sp.BaggageItem("error4")

	// sp2 := opentracing.StartSpan("Meio1", opentracing.ChildOf(sp.Context()))
	// time.Sleep(1 * time.Second)
	// defer sp2.Finish()

	// sp21 := opentracing.StartSpan("Meio2", opentracing.ChildOf(sp.Context()))
	// time.Sleep(1 * time.Second)
	// defer sp21.Finish()

	// sp22 := opentracing.StartSpan("Meio3", opentracing.ChildOf(sp.Context()))
	// time.Sleep(1 * time.Second)
	// defer sp22.Finish()

	// sp23 := opentracing.StartSpan("Meio4", opentracing.ChildOf(sp.Context()))
	// time.Sleep(1 * time.Second)
	// defer sp23.Finish()

	// sp3, ctx := opentracing.StartSpanFromContext(ctx, "Fim")
	// time.Sleep(1 * time.Second)
	// defer sp3.Finish()

	return c.SendStatus(fiber.StatusNotFound)
}
