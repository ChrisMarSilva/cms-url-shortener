package main

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// go test
// go test -v
// go test -bench=.
// go test -bench=Add
// go test -run TestUserrepoGetAll -v
// go test -run=XXX -bench . -benchmem

func TestRouteShortenStatus200(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg": "ok"})
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/test", nil))
	// req := httptest.NewRequest("GET", "/test", nil)
	// resp, err := app.Test(req, 1)

	utils.AssertEqual(t, nil, err, "app.Test")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
}
