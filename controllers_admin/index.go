package controllers_admin

import (
	"context"
	"log"
	"vauntly/components"
	utils "vauntly/utils"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/iterator"

	"github.com/labstack/echo/v4"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func Home(c echo.Context) error {
	config, error := utils.GetConfig(c)

	if error != nil {
		return error
	}

	app := c.Get("firebase").(*firebase.App)
	if app == nil {
		log.Fatalf("error getting firebase app from context\n")
	}

	ctx := context.Background()
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	iter := client.Users(ctx, "")
	users := []auth.ExportedUserRecord{}
	for {
		user, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// log.Fatalf("error listing users: %s\n", err)
			return components.Layout("Admin", config,
				Div(
					H1(g.Text("Admin")),
					Div(
						P(
							A(
								g.Attr("href", "/admin"),
								g.Text("Reload"),
							),
						),
						P(
							g.Text("Error listing users"),
						),
						P(
							g.Text(err.Error()),
						),
					),
				)).Render(c.Response().Writer)
		}
		users = append(users, *user)
		log.Printf("read user user: %v\n", user.Email)
	}

	return components.Layout("Admin", config,
		Div(
			P(
				A(
					g.Attr("href", "/admin"),
					g.Text("Reload"),
				),
			),
			Div(
				Class("grid grid-cols-3 gap-2 bg-black text-white text-sm p-4 rounded-lg"),
				Span(
					Class("text-red-400"),
					g.Text("Created On")),
				Span(
					Class("text-red-400"),
					g.Text("Email")),
				Span(

					Class("text-red-400"),
					g.Text("Display Name")),
				g.Group(g.Map(users, func(user auth.ExportedUserRecord) g.Node {
					return g.Group(
						g.Map([]string{"CreatedOn", "Email", "DisplayName"}, func(key string) g.Node {
							if key == "CreatedOn" {
								return Span(g.Textf("%d", user.UserMetadata.CreationTimestamp))
							}
							if key == "Email" {
								return Span(
									Class("font-semibold"),
									g.Text(user.Email),
								)
							}
							if key == "DisplayName" {
								return Span(g.Text(user.DisplayName))
							}
							return nil
						}),
					)
				})),
			),
		),
	).Render(c.Response().Writer)
}
