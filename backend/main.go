package main

import (
	"backend/db"
	"backend/routes"
	"context"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	client := db.NewClient()
	ctx := context.Background()

	client.Prisma.Connect()
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	// JWT Middleware
	middleware := (jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	}))

	api := app.Group("/api")
	routes.AuthRouter(api, client, ctx, middleware)
	routes.UserRouter(api, client, ctx, middleware)
	routes.TweetRouter(app, client, ctx, middleware)

	app.Listen(":3000")
}
