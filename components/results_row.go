package components

import (
	"vauntly/models"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

func ResultsRow(post models.Post) g.Node {
	class := "search_row w-full group relative mb-2 rounded-sm border-0 border-pink-200 bg-pink-100 dark:border-prussian-blue-900 dark:bg-slate-700 dark:text-white"
	if post.IsPinned() {
		class = "search_row group relative mb-2 rounded-sm border-0 border-orange-200 bg-orange-200 text-black dark:border-slate-700 dark:bg-slate-900 dark:text-white"
	}

	return h.Div(
		h.Class(class),
		g.Attr("onclick", "utils.toggleOpenState('cbx"+post.ID+"', 'desc"+post.ID+"')"),
		h.Div(
			c.Classes{"cursor-pointer flex h-32 items-center space-x-4 px-4": true},
			g.If(post.Thumbnail != "", h.Span(
				c.Classes{"rounded-full initials inline-flex h-[40px] w-[40px] my-2 shrink-0 items-center justify-center overflow-hidden": true},
				h.Img(h.Src(post.Thumbnail), h.Width("40"), h.Height("40")),
			)),
			g.If(post.Thumbnail == "", h.Span(
				c.Classes{"bg-teal-300 rounded-full initials inline-flex h-[40px] w-[40px] my-2 shrink-0 items-center justify-center overflow-hidden": true},
				g.Text(post.GetInitials()),
			)),
			h.Div(
				c.Classes{"flex grow": true},
				h.Div(
					c.Classes{"flex grow flex-col min-w-10": true},

					h.P(
						c.Classes{"text-black line-clamp-1 font-semibold lg:line-clamp-2": true},
						g.Text(post.CompanyName),
					),
					h.P(
						c.Classes{"text-black line-clamp-1 font-bold md:line-clamp-2": true},
						g.Text(post.Title),
					),
					h.P(
						c.Classes{"text-black": true},
						g.Text(post.Location),
					),
				),
				h.Div(
					h.Class("tag-container"),
					g.Group(g.Map(post.Tags, func(tag string) g.Node {
						return h.Button(
							g.Attr("onclick", "utils.halt(event)"),
							h.Class("hidden cursor-pointer rounded-md bg-white px-2 mr-1 font-semibold text-pink-950 transition-colors duration-300 hover:bg-blue-100 hover:text-black my-[2px]"),
							g.Text(tag),
						)
					})),
				),
			),
			h.Span(
				h.Class("m-2"),
				g.Text(post.TimeSinceCreated()),
			),
			h.Span(
				g.Attr("onclick", "utils.halt(event)"),
				h.Button(
					h.Class("btn-apply done"),
					g.Text("Apply"),
				),
			),
		),
		h.Div(
			c.Classes{"flex flex-col items-center justify-center": true},
			h.Input(
				h.Type("checkbox"),
				h.ID("cbx"+post.ID),
				g.Attr("aria-label", "toggle show description"),
				h.Class("peer hidden"),
			),
			h.Div(
				h.Class("hidden p-4 peer-checked:flex"),
				h.Div(
					h.Class("items-center justify-center"),
					g.Attr("hx-get", "/posts/details/"+post.ID+"/"),
					g.Attr("hx-indicator", "#htmx"+post.ID),
					g.Attr("hx-swap", "outerHTML transition:true"),
					g.Attr("hx-trigger", "change"),
					h.ID("desc"+post.ID),
					h.Img(
						h.ID("htmx"+post.ID),
						h.Src("/public/images/bars-loader.svg"),
						h.Height("48"),
						h.Width("48"),
					),
				),
			),
		),
	)
}
