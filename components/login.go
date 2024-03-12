package components

import (
	"vauntly/models"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

func Login(config *models.Config) g.Node {
	return Layout(
		"Login",
		config,
		h.Div(
			h.Class("flex flex-col items-center justify-center"),
			h.Div(
				h.Class("w-96"),
				h.FormEl(
					h.Class("flex flex-col gap-4"),
					h.Input(
						h.Type("text"),
						h.Name("email"),
						h.Placeholder("Email"),
						h.Class("px-4 py-2 border border-gray-300 rounded-md"),
					),
					h.Input(
						h.Type("password"),
						h.Name("password"),
						h.Placeholder("Password"),
						h.Class("px-4 py-2 border border-gray-300 rounded-md"),
					),
					h.Button(
						h.Type("submit"),
						h.Class("bg-primary-700 text-white px-4 py-2 rounded-md"),
						g.Text("Log in"),
					),
				),
			),
		),
	)
}
