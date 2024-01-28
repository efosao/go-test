package middleware

import (
	"gofiber/models"

	"github.com/gofiber/fiber"
)

func SetupThemes(c *fiber.Ctx) error {
	if isGet := (c.Method() == "GET"); isGet {
		theme := c.Cookies("theme")
		themeOptions := []models.ThemeOption{}
		themeOptions = append(themeOptions, models.ThemeOption{Value: "light", Label: "🌞", Selected: theme == "light"})
		themeOptions = append(themeOptions, models.ThemeOption{Value: "dark", Label: "🌘", Selected: theme == "dark"})
		themeOptions = append(themeOptions, models.ThemeOption{
			Value:    "system",
			Label:    "🌎",
			Selected: (theme != "light" && theme != "dark"),
		})
		c.Locals("ThemeOptions", themeOptions)
	}
	c.Next()
	return nil
}
