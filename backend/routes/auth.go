package routes

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"backend/db"
)

func accessible(c *fiber.Ctx) error {
	return c.SendString("Accessible")
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}

func login(client *db.PrismaClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		handle := c.FormValue("Handle")
		password := c.FormValue("password")

		user, uerr := client.User.FindUnique(
			db.User.Handle.Equals(handle),
		).Exec(ctx)

		if uerr != nil {
			return c.SendStatus(fiber.ErrBadRequest.Code)
		}

		if user.Password != password {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		claims := jwt.MapClaims{
			"name": user.Handle,
			"exp":  time.Now().Add(time.Hour * 1000).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// curl -X POST -d "handle=John&password=doe" http://localhost:3000/login
		return c.JSON(fiber.Map{"token": t})
	}
}

func register(client *db.PrismaClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		handle := c.FormValue("handle")
		password := c.FormValue("password")
		displayName := c.FormValue("displayName")

		_, uerr := client.User.CreateOne(
			db.User.Handle.Set(handle),
			db.User.DisplayName.Set(displayName),
			db.User.Password.Equals(password),
		).Exec(ctx)

		if uerr != nil {
			return c.SendStatus(fiber.ErrInternalServerError.Code)
		}

		return c.SendStatus(200)
	}
}

func AuthRouter(app fiber.Router, client *db.PrismaClient, ctx context.Context, middlware fiber.Handler) {
	auth := app.Group("/auth")
	auth.Post("/login", login(client, ctx))
	auth.Post("/register", register(client, ctx))
	auth.Get("/accessible", accessible)
	auth.Use(middlware)
	auth.Get("/restricted", restricted)
}
