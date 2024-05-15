package main

import (
	_ "embed"
	c "vauntly/controllers"
	ca "vauntly/controllers_admin"
	mw "vauntly/middleware"
	"vauntly/models"
	"vauntly/utils"

	"context"
	"fmt"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
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

func initFirebase() (*firebase.App, error) {
	firebaseServiceAccount := os.Getenv("FIREBASE_SERVICE_ACCOUNT")
	opt := option.WithCredentialsJSON([]byte(firebaseServiceAccount))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	return app, err
}

func main() {
	models.ConnectDB()
	e := echo.New()

	// store the firebase app in the echo context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			app, err := initFirebase()
			if err != nil {
				return err
			}
			c.Set("firebase", app)
			return next(c)
		}
	})

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
	e.POST("/login", c.LoginPost)
	e.GET("/admin", ca.Home)

	PORT := ":8000"
	if os.Getenv("PORT") != "" {
		PORT = fmt.Sprintf(":%s", os.Getenv("PORT"))
	}

	println("Server running on port", PORT)
	e.Logger.Fatal(e.Start(PORT))
}
