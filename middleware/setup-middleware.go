package middleware

import (
	"log"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupMiddleware(e *echo.Echo) {
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
}
