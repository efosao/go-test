package middleware

import (
	"context"
	"gofiber/models"
	"net/http"
)

func UserTheme(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		for _, cookie := range r.Cookies() {
			theme := cookie.Value
			switch cookie.Name {
			case "theme":
				themeOptions := []models.ThemeOption{}
				themeOptions = append(themeOptions, models.ThemeOption{Value: "light", Label: "🌞", Selected: theme == "light"})
				themeOptions = append(themeOptions, models.ThemeOption{Value: "dark", Label: "🌘", Selected: theme == "dark"})
				themeOptions = append(themeOptions, models.ThemeOption{
					Value:    "system",
					Label:    "🌎",
					Selected: (theme != "light" && theme != "dark"),
				})
				ctx = context.WithValue(ctx, models.ThemeOptionsKey, themeOptions)
				ctx = context.WithValue(ctx, models.ThemeKey, theme)
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
