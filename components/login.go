package components

import (
	"vauntly/models"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

func Login(config *models.Config, props models.LoginProps) g.Node {
	return Layout(
		"Login",
		config,
		h.Div(
			h.Class("flex flex-col items-center justify-center flex-grow h-[calc(100vh-8rem)]"),
			h.Div(
				h.Class("w-96 rounded-md bg-slate-500 p-4 form"),
				h.FormEl(
					h.Method("POST"),
					g.Attr("up-submit"),
					g.Attr("up-target", ":none"),
					h.Action("/login"),
					h.Class("flex flex-col gap-4"),
					h.H1(
						h.Class("text-3xl font-bold text-white"),
						g.Text("Log in"),
					),
					h.Input(
						h.Type("text"),
						h.Name("email"),
						h.AutoComplete("email"),
						h.Placeholder("Email"),
						g.Attr("up-validate"),
						h.Value(props.Email),
						h.Class("px-4 py-2 border border-gray-300 rounded-md"),
					),
					g.If(props.EmailErr != "", h.Div(
						h.Class("text-red-700 font-bold px-2"),
						g.Text(props.EmailErr),
					)),
					h.Input(
						h.Type("password"),
						h.Name("password"),
						g.Attr("up-validate"),
						h.Value(props.Password),
						h.AutoComplete("current-password"),
						h.Placeholder("Password"),
						h.Class("px-4 py-2 border border-gray-300 rounded-md"),
					),
					h.Button(
						h.Type("submit"),
						h.Class("bg-purple-700 text-white px-4 py-2 rounded-md"),
						g.Text("Log in"),
					),
				),
			),
		),
	)
}
