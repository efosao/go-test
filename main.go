package main

import (
	"gofiber/controllers"
	"gofiber/middleware"
	"gofiber/models"

	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/mustache/v2"
)

func main() {
	models.ConnectDB()
	engine := mustache.New("./views", ".mustache")
	app := fiber.New(fiber.Config{Views: engine})

	app.Use(compress.New())
	app.Use(helmet.New(
		helmet.Config{
			ContentSecurityPolicy: `
				default-src 'self';
				base-uri 'self';
				font-src 'self' https: data:;
				form-action 'self';
				frame-ancestors 'self';
				img-src 'self' https: data:;
				object-src 'none';
				script-src 'self' unpkg.com cdn.jsdelivr.net 'unsafe-eval' 'unsafe-inline';
				script-src-attr 'unsafe-inline';
				style-src 'self' https: 'unsafe-inline';
				upgrade-insecure-requests
			`,
		}))
	app.Use(etag.New())
	app.Use(recover.New())
	app.Use(middleware.SetupThemes)

	app.Static("/public", "./public", fiber.Static{
		Compress:      true,
		CacheDuration: 60 * time.Second,
		ByteRange:     true,
		ModifyResponse: func(ctx *fiber.Ctx) error {
			ctx.Set(fiber.HeaderCacheControl, fmt.Sprintf("private, max-age=%d, stale-while-revalidate=3600", 3600))
			return nil
		},
	})

	app.Get("/", controllers.GetHome)
	app.Get("/monitor", monitor.New())
	app.Get("/posts", controllers.GetPosts)

	partials := app.Group("/partials")
	partials.Get("/posts/details/:id", controllers.GetPostDetail)
	partials.Post("/posts/search/:page", controllers.PostSearchResultsPage)

	PORT := "localhost:8000"
	if os.Getenv("PORT") != "" {
		PORT = fmt.Sprintf(":%s", os.Getenv("PORT"))
	}

	log.Fatal(app.Listen(PORT))
}
