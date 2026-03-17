package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	// A simple health check or welcome route
	app.Get("/", func(c fiber.Ctx) error {
		c.SendString("Todo API is running! great new ")
		return c.SendStatus(200)
	})

app.Get("/api/v1", func(c fiber.Ctx) error {
    return c.Status(200).JSON(fiber.Map{
        "msg": "hello world",
    })
})

	// Example: Get all todos (placeholder)
	// c -> context
	app.Get("/api/todos", func(c fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"todos": []string{}})
	})

	log.Fatal(app.Listen(":4000"))
}
