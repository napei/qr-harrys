package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
)

// TODO Implement a GOSIP client https://go.spflow.com/samples/library-initiation

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(requestid.New())

	app.Get("/:id", func(c *fiber.Ctx) error {
		url := "https://google.com.au"
		return c.Redirect(url, fiber.StatusMovedPermanently)
	})

	app.Listen(":3000")
}
