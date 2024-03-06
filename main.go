package main

import (
	_ "embed"
	c "vauntly/controllers"
	mw "vauntly/middleware"
	"vauntly/models"

	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

//go:generate sh -c "printf %s $(git rev-parse HEAD) > hash.txt"
//go:embed hash.txt
var cacheHash string

func main() {
	models.ConnectDB()

	r := http.NewServeMux()
	fs := http.FileServer(http.Dir("public"))
	c.SetupCacheHash(cacheHash)

	r.Handle("GET /public/", http.StripPrefix("/public/", fs))
	r.Handle("GET /{$}", mw.UserTheme(http.HandlerFunc(c.GetHome)))
	r.Handle("GET /about/{$}", mw.UserTheme(http.HandlerFunc(c.GetAbout)))
	r.Handle("POST /about/{$}", mw.UserTheme(http.HandlerFunc(c.PostAbout)))
	r.Handle("GET /posts/{$}", mw.UserTheme(http.HandlerFunc(c.GetPosts)))
	r.Handle("GET /posts/details/{id}", mw.UserTheme(http.HandlerFunc(c.GetPostDetail)))
	r.Handle("GET /partials/posts/search/{page}", mw.UserTheme(http.HandlerFunc(c.PostSearchResultsPage)))
	r.Handle("POST /partials/posts/search/{page}", mw.UserTheme(http.HandlerFunc(c.PostSearchResultsPage)))

	PORT := ":8000"
	if os.Getenv("PORT") != "" {
		PORT = fmt.Sprintf(":%s", os.Getenv("PORT"))
	}

	println("Server running on port", PORT)
	log.Fatal(http.ListenAndServe(PORT, handlers.CompressHandler(r)))
}
