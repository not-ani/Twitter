package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"backend/db"
)

func h(client *db.PrismaClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	}
}

func TodoRouter(app fiber.Router, client *db.PrismaClient, ctx context.Context, middlware fiber.Handler) {
	todo := app.Group("/todo")

	todo.Get("/", h(client, ctx))
}
