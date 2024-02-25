package middleware

import (
	"context"
	"gofiber/models"
	"net/http"
)

func ReadThemeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		cookiesMap := map[string]string{}
		for _, cookie := range cookies {
			cookiesMap[cookie.Name] = cookie.Value
		}
		theme := cookiesMap["theme"]
		themeOptions := []models.ThemeOption{}
		themeOptions = append(themeOptions, models.ThemeOption{Value: "light", Label: "🌞", Selected: theme == "light"})
		themeOptions = append(themeOptions, models.ThemeOption{Value: "dark", Label: "🌘", Selected: theme == "dark"})
		themeOptions = append(themeOptions, models.ThemeOption{
			Value:    "system",
			Label:    "🌎",
			Selected: (theme != "light" && theme != "dark"),
		})

		ctx := context.WithValue(r.Context(), models.ThemeOptionsKey, themeOptions)
		ctx = context.WithValue(ctx, models.ThemeKey, theme)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
