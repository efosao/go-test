package utils

import (
	"vauntly/models"

	"github.com/labstack/echo/v4"
)

var CacheHash = ""

func SetupCacheHash(hashValue string) {
	if len(hashValue) < 8 {
		CacheHash = hashValue
		return
	}
	CacheHash = hashValue[0:8]
}

func GetConfig(c echo.Context) (*models.Config, error) {
	themeCookie, _ := c.Cookie("theme")
	theme := "system"
	if themeCookie != nil {
		theme = themeCookie.Value
	}

	themeOptions := []models.ThemeOption{
		{Value: "light", Label: "🌞", Selected: theme == "light"},
		{Value: "dark", Label: "🌘", Selected: theme == "dark"},
		{Value: "system", Label: "🌎", Selected: theme != "light" && theme != "dark"},
	}

	config := &models.Config{
		Path:         c.Request().URL.Path,
		Theme:        theme,
		ThemeOptions: themeOptions,
	}
	return config, nil
}
