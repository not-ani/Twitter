package routes

import (
	"backend/db"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type CreateTweetParams struct {
	Content string `json:"content"`
}

type FindOneTweetParams struct {
	ID string `json:"id"`
}

type UpdateTweetParams struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func create(client *db.PrismaClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var tweet CreateTweetParams

		if err := c.BodyParser(&tweet); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		name := claims["name"].(string)

		_, err := client.Tweets.CreateOne(
			db.Tweets.Content.Set(tweet.Content),
			db.Tweets.Author.Link(db.User.Handle.Equals(name)),
		).Exec(ctx)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendStatus(200)
	}
}

func findOne(client *db.PrismaClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		tweet, userr := client.Tweets.FindUnique(
			db.Tweets.ID.Equals(id),
		).Exec(ctx)

		if userr != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		return c.JSON(tweet)
	}
}

func update(client *db.PrismaClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var tweet UpdateTweetParams
		if err := c.BodyParser(&tweet); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		_, err := client.Tweets.FindUnique(
			db.Tweets.ID.Equals(tweet.ID),
		).Update(
			db.Tweets.Content.Set(tweet.Content),
		).Exec(ctx)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendStatus(200)
	}
}

func findAll(client *db.PrismaClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		tweets, err := client.Tweets.FindMany().Exec(ctx)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(tweets)
	}
	// curl http://localhost:3000/tweet/all
}

func delete(client *db.PrismaClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		_, err := client.Tweets.FindUnique(
			db.Tweets.ID.Equals(id),
		).Delete().Exec(ctx)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendStatus(200)
	}
}

func TweetRouter(app fiber.Router, client *db.PrismaClient, ctx context.Context, middlware fiber.Handler) {
	client.Prisma.Connect()
	tweet := app.Group("/tweet")
	tweet.Get("/all", findAll(client, ctx))
	tweet.Get("/:id", findOne(client, ctx))
	tweet.Use(middlware)
	tweet.Patch("/", update(client, ctx))
	tweet.Post("/", create(client, ctx))
	tweet.Delete("/:id", delete(client, ctx))
}
