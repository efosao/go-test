package main

import (
	_ "embed"
	c "vauntly/controllers"
	"vauntly/models"

	"net/http"

	"github.com/labstack/echo/v4"
)

//go:generate sh -c "printf %s $(git rev-parse HEAD) > hash.txt"
//go:embed hash.txt
var cacheHash string

func main() {
	models.ConnectDB()

	e := echo.New()
	fs := http.FileServer(http.Dir("public"))
	c.SetupCacheHash(cacheHash)

	e.GET("/public/*", echo.WrapHandler(http.StripPrefix("/public/", fs)))
	e.GET("/", c.GetHome)
	e.GET("/about", c.GetAbout)
	e.POST("/about", c.PostAbout)
	e.GET("/posts", c.GetPosts)
	e.GET("/posts/details/:id", c.GetPostDetail)
	e.GET("/partials/posts/search/:page", c.PostSearchResultsPage)
	e.POST("/partials/posts/search/:page", c.PostSearchResultsPage)

	PORT := ":8000"
	println("Server running on port", PORT)
	e.Logger.Fatal(e.Start(":3000"))
}
