package middleware

import (
	"context"
	"net/http"
	"vauntly/models"
)

func UserTheme(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		for _, cookie := range r.Cookies() {
			switch cookie.Name {
			case "theme":
				themeOptions := []models.ThemeOption{}
				theme := ""
				if cookie.Value == "light" || cookie.Value == "dark" {
					theme = cookie.Value
				}

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

		if ctx.Value(models.ThemeKey) == nil {
			themeOptions := []models.ThemeOption{
				{Value: "light", Label: "🌞", Selected: false},
				{Value: "dark", Label: "🌘", Selected: false},
				{Value: "system", Label: "🌎", Selected: true},
			}

			ctx = context.WithValue(ctx, models.ThemeOptionsKey, themeOptions)
			ctx = context.WithValue(ctx, models.ThemeKey, "")
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
