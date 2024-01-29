package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"backend/db"
)

type FindFollwingParams struct {
	UserHandle string `json:"id"`
}

func findAllFollowingOfUser(client *db.PrismaClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var params FindFollwingParams
		if err := c.BodyParser(&params); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		// Query the database to find the following users
		following, err := client.FollowRelation.FindMany(
			db.FollowRelation.FollowerHandle.Equals(params.UserHandle),
		).Exec(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"following": following,
		})
	}
}

func findAllFollwersOfUser(client db.PrismaClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var params FindFollwingParams
		if err := c.BodyParser(&params); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		// Query the database to find the following users
		following, err := client.FollowRelation.FindMany(
			db.FollowRelation.FollowingHandle.Equals(params.UserHandle),
		).Exec(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"following": following,
		})
	}
}

func UserRouter(app fiber.Router, client *db.PrismaClient, ctx context.Context, middlware fiber.Handler) {
	users := app.Group("/users")

	users.Post("/following", middlware, findAllFollowingOfUser(client, ctx))
	users.Post("/followers", middlware, findAllFollwersOfUser(*client, ctx))
}
