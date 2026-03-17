package main

import (
	"fmt"
	"log"
	"os"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	app := fiber.New()
	todos := []Todo{}
	// load env vars
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading env Vars")
	}
	PORT := os.Getenv("PORT")

	// create TODO route
	app.Post("/api/v1/todos", func(c fiber.Ctx) error {
		todo := Todo{}
		if err := c.Bind().Body(&todo); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}
		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "body is required"})
		}
		todo.ID = len(todos) + 1
		todos = append(todos, todo)
		return c.Status(201).JSON(todo)
	})
	// Update TODO
	app.Patch("/api/v1/todos/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = !todos[i].Completed
				return c.Status(200).JSON(todos[i])
			}
		}
		return c.Status(404).JSON(fiber.Map{"msg": "Todo not found"})
	})
	// Delete TODO
	app.Delete("/api/v1/todos/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				// what happen here he take todos after and before and and make sperate of it
				// veritake one
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"msg": "Todo deleted"})
			}
		}
		return c.Status(404).JSON(fiber.Map{"msg": "Todo not found"})
	})
	// Get all todos
	app.Get("/api/v1/todos", func(c fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"todos": todos,
		})
	})

	// app.Get("/api/todos", func(c fiber.Ctx) error {
	// 		return c.Status(200).JSON(fiber.Map{"todos": []string{}})
	// })

	app.Get("/api/v1", func(c fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"msg": "hello world",
		})
	})

	log.Fatal(app.Listen(":" + PORT))

}
