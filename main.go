package main

import (
	c "vauntly/controllers"
	mw "vauntly/middleware"
	"vauntly/models"
	"vauntly/utils"

	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
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

	mw.SetupMiddleware(e)
	fs := http.FileServer(http.Dir("public"))
	utils.SetupCacheHash(cacheHash)
	e.HTTPErrorHandler = customHTTPErrorHandler

	e.GET("/public/*", echo.WrapHandler(http.StripPrefix("/public/", fs)))
	e.GET("/", c.Home)
	e.GET("/about", c.GetAbout)
	e.POST("/about", c.PostAbout)
	e.GET("/posts", c.GetPosts)
	e.GET("/posts/details/:id", c.GetPostDetail)
	e.GET("/partials/posts/search/:page", c.PostSearchResultsPage)
	e.POST("/partials/posts/search/:page", c.PostSearchResultsPage)
	e.GET("/login", c.Login)

	PORT := ":8000"
	if os.Getenv("PORT") != "" {
		PORT = fmt.Sprintf(":%s", os.Getenv("PORT"))
	}

	println("Server running on port", PORT)
	e.Logger.Fatal(e.Start(PORT))
}
