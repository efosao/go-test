package main

import (
	_ "embed"
	"log"
	"strings"
	"time"
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

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
}

func main() {
	models.ConnectDB()

	e := echo.New()
	e.Pre(middleware.RewriteWithConfig(middleware.RewriteConfig{
		Rules: map[string]string{
			"*/":   "$1",
			"*/?*": "$1?$2",
		},
	}))
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "custom timeout error message returns to client",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			log.Println(c.Path())
		},
		Timeout: 20 * time.Second,
	}))
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

	e.HTTPErrorHandler = customHTTPErrorHandler

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
