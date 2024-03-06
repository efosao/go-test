package main

import (
	_ "embed"
	"strings"
	c "vauntly/controllers"
	"vauntly/models"

	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:generate sh -c "printf %s $(git rev-parse HEAD) > hash.txt"
//go:embed hash.txt
var cacheHash string

func main() {
	models.ConnectDB()

	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${status} ${method} ${uri} ${latency_human}\n",
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "/public")
		},
	}))
	e.Use(middleware.Recover())

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
	if os.Getenv("PORT") != "" {
		PORT = fmt.Sprintf(":%s", os.Getenv("PORT"))
	}

	println("Server running on port", PORT)
	e.Logger.Fatal(e.Start(PORT))
}
