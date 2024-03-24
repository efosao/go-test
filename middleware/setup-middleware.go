package middleware

import (
	"log"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupMiddleware sets up the middleware for the given Echo instance.
// It configures various middleware functions such as Rewrite, Timeout, Gzip, Logger, and Recover.
// The Rewrite middleware rewrites the URL path based on the specified rules.
// The Timeout middleware sets a timeout for each request and handles timeout errors.
// The Gzip middleware compresses the response body using gzip compression.
// The Logger middleware logs each request with the specified format.
// The Recover middleware recovers from panics and returns an HTTP 500 error.
func SetupMiddleware(e *echo.Echo) {
	e.Pre(middleware.RewriteWithConfig(middleware.RewriteConfig{
		Rules: map[string]string{
			"*/":   "$1",
			"*/?*": "$1?$2",
		},
	}))
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Server request timed out\n",
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
