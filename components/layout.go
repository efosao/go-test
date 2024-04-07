package components

import (
	"fmt"
	"vauntly/models"
	"vauntly/utils"

	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

func Layout(title string, config *models.Config, children g.Node) g.Node {
	releaseHash := fmt.Sprintf("?v=%s", utils.CacheHash)
	return h.Doctype(
		h.HTML(
			h.Class("smooth-scroll "+config.Theme),
			h.Lang("en"),
			h.Head(
				h.TitleEl(g.Text(title)),
				h.StyleEl(h.Type("text/css"), g.Raw(".is-active{ font-weight: bold }")),
				h.Link(h.Rel("stylesheet"), h.Href(fmt.Sprintf("/public/dist/stylesheet.css%s", releaseHash))),
				h.Script(h.Src(fmt.Sprintf("/public/dist/index.js%s", releaseHash)), h.Defer()),
				h.Meta(h.Name("viewport"), h.Content("width=device-width, initial-scale=1")),
				g.Raw(`<link rel="apple-touch-icon" sizes="180x180" href="/public/apple-touch-icon.png">`),
				g.Raw(`<link rel="icon" type="image/png" sizes="32x32" href="/public/favicon-32x32.png">`),
				g.Raw(`<link rel="icon" type="image/png" sizes="16x16" href="/public/favicon-16x16.png">`),
				h.Link(h.Rel("manifest"), h.Href("/public/site.webmanifest")),
			),
			h.Body(
				h.Class("flex flex-col min-h-screen dark:bg-slate-400"),
				hx.Boost("true"),
				g.If(config.ShowNav, Navbar(config)),
				h.Div(
					h.Class("bg-orange-300 transition-colors dark:bg-slate-600 h-2"),
					h.Div(
						h.Class("mx-auto max-w-5xl"),
						// g.Raw("<app-bar><div class='h-10 bg-black rounded-md'></div></app-bar>"),
					),
				),
				h.Div(
					h.Class("flex-grow w-full p-2 mx-auto max-w-5xl"),
					h.ID("page-content"),
					g.If(config.ShowNav,
						h.H1(
							h.ID("page-title"),
							h.Class("my-4 overflow-hidden max-w-5xl text-3xl font-extrabold text-black dark:text-black"),
							g.Text(title),
						)),
					children,
				),
				h.Footer(
					h.Class("flex justify-center gap-2 mt-4 mb-2"),
					h.Span(
						h.Class("text-md font-bold text-gray-900 bg-slate-300 p-2 rounded-md"),
						g.Text("Release:"),
						h.Span(
							h.Class("text-md font-bold text-gray-700 p-1 ml-2 bg-slate-400 rounded-md"),
							g.Text(utils.CacheHash),
						),
					),
				),
			),
		),
	)
}
