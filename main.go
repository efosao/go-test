package main

import (
	"gofiber/controllers"
	"gofiber/models"

	mw "gofiber/middleware"

	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func main() {
	models.ConnectDB()

	r := http.NewServeMux()
	fs := http.FileServer(http.Dir("public"))

	r.Handle("GET /public/", http.StripPrefix("/public/", fs))
	r.Handle("GET /{$}", mw.ReadThemeMiddleware(http.HandlerFunc(controllers.GetHome)))
	r.Handle("GET /about/{$}", mw.ReadThemeMiddleware(http.HandlerFunc(controllers.GetAbout)))
	r.Handle("GET /posts/{$}", mw.ReadThemeMiddleware(http.HandlerFunc(controllers.GetPosts)))
	r.Handle("GET /posts/details/{id}", mw.ReadThemeMiddleware(http.HandlerFunc(controllers.GetPostDetail)))
	r.Handle("GET /partials/posts/search/{page}", mw.ReadThemeMiddleware(http.HandlerFunc(controllers.PostSearchResultsPage)))
	r.Handle("POST /partials/posts/search/{page}", mw.ReadThemeMiddleware(http.HandlerFunc(controllers.PostSearchResultsPage)))

	PORT := ":8000"
	if os.Getenv("PORT") != "" {
		PORT = fmt.Sprintf(":%s", os.Getenv("PORT"))
	}

	println("Server running on port", PORT)

	log.Fatal(http.ListenAndServe(PORT, handlers.CompressHandler(r)))
}
