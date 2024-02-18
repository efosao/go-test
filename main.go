package main

import (
	"context"
	"gofiber/controllers"
	"gofiber/middleware"
	"gofiber/models"

	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing middlewareOne")

		// theme := c.Cookies("theme")
		cookies := r.Cookies()
		cookiesMap := map[string]string{}
		for _, cookie := range cookies {
			cookiesMap[cookie.Name] = cookie.Value
		}
		theme := cookiesMap["theme"]
		themeOptions := []models.ThemeOption{}
		themeOptions = append(themeOptions, models.ThemeOption{Value: "light", Label: "🌞", Selected: theme == "light"})
		themeOptions = append(themeOptions, models.ThemeOption{Value: "dark", Label: "🌘", Selected: theme == "dark"})
		themeOptions = append(themeOptions, models.ThemeOption{
			Value:    "system",
			Label:    "🌎",
			Selected: (theme != "light" && theme != "dark"),
		})

		ctx := context.WithValue(r.Context(), "themeOption", themeOptions)
		ctx = context.WithValue(ctx, "theme", theme)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	models.ConnectDB()
	engine := html.New("./views", ".tpl")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(compress.New())
	app.Use(logger.New())
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

	fs := http.FileServer(http.Dir("public"))
	http.Handle("GET /{$}", middlewareOne(http.HandlerFunc(controllers.GetAbout)))
	http.Handle("GET /about/{$}", middlewareOne(http.HandlerFunc(controllers.GetAbout)))
	http.Handle("GET /posts/{$}", middlewareOne(http.HandlerFunc(controllers.GetAbout)))
	http.Handle("GET /public/", http.StripPrefix("/public/", fs))

	app.Get("/", controllers.GetHome)
	app.Post("/", controllers.GetHome)
	app.Get("/posts", controllers.GetPosts)
	app.Get("/monitor", monitor.New())

	partials := app.Group("/partials")
	partials.Get("/posts/details/:id", controllers.GetPostDetail)
	partials.Post("/posts/search/:page", controllers.PostSearchResultsPage)

	PORT := "localhost:8000"
	if os.Getenv("PORT") != "" {
		PORT = fmt.Sprintf(":%s", os.Getenv("PORT"))
	}

	log.Fatal(http.ListenAndServe("localhost:8080", nil))

	log.Fatal(app.Listen(PORT))
}
