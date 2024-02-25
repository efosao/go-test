package main

import (
	"context"
	"gofiber/controllers"
	"gofiber/models"

	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func readThemeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx := context.WithValue(r.Context(), models.ThemeOptionsKey, themeOptions)
		ctx = context.WithValue(ctx, models.ThemeKey, theme)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	models.ConnectDB()

	r := http.NewServeMux()
	fs := http.FileServer(http.Dir("public"))

	r.Handle("GET /public/", http.StripPrefix("/public/", fs))
	r.Handle("GET /{$}", readThemeMiddleware(http.HandlerFunc(controllers.GetHome)))
	r.Handle("GET /about/{$}", readThemeMiddleware(http.HandlerFunc(controllers.GetAbout)))
	r.Handle("GET /posts/{$}", readThemeMiddleware(http.HandlerFunc(controllers.GetPosts)))
	r.Handle("GET /posts/details/{id}", readThemeMiddleware(http.HandlerFunc(controllers.GetPostDetail)))
	r.Handle("GET /partials/posts/search/{page}", readThemeMiddleware(http.HandlerFunc(controllers.PostSearchResultsPage)))
	r.Handle("POST /partials/posts/search/{page}", readThemeMiddleware(http.HandlerFunc(controllers.PostSearchResultsPage)))

	PORT := ":8000"
	if os.Getenv("PORT") != "" {
		PORT = fmt.Sprintf(":%s", os.Getenv("PORT"))
	}

	println("Server running on port", PORT)

	log.Fatal(http.ListenAndServe(PORT, handlers.CompressHandler(r)))
}
