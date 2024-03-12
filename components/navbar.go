package components

import (
	"vauntly/models"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

func Navbar(config *models.Config) g.Node {
	currentPath := config.Path
	return h.Div(
		h.Nav(
			h.Class("bg-white border-gray-200 px-4 lg:px-6 py-2.5 dark:bg-gray-800"),
			h.Div(
				h.Class("flex flex-wrap justify-between items-center mx-auto max-w-screen-xl"),
				h.A(
					h.Href("/"),
					h.Class("flex items-center"),
					h.Img(
						h.Src("/public/test-logo.svg"),
						h.Class("mr-3 h-6 sm:h-9"),
						h.Alt("Flowbite Logo"),
					),
					h.Span(
						h.Class("self-center text-xl font-semibold whitespace-nowrap dark:text-white"),
						g.Text("Vauntly"),
					),
				),
				h.Div(
					h.Class("flex items-center lg:order-2"),
					h.A(
						h.Href("#"),
						h.Class("text-gray-800 dark:text-white hover:bg-gray-50 focus:ring-4 focus:ring-gray-300 font-medium rounded-lg text-sm px-4 lg:px-5 py-2 lg:py-2.5 mr-2 dark:hover:bg-gray-700 focus:outline-none dark:focus:ring-gray-800"),
						g.Text("Log in"),
					),
					h.A(
						h.Href("#"),
						h.Class("text-white bg-primary-700 hover:bg-primary-800 focus:ring-4 focus:ring-primary-300 font-medium rounded-lg text-sm px-4 lg:px-5 py-2 lg:py-2.5 mr-2 dark:bg-primary-600 dark:hover:bg-primary-700 focus:outline-none dark:focus:ring-primary-800"),
						g.Text("Get started"),
					),
					h.Button(
						h.DataAttr("collapse-toggle", "mobile-menu-2"),
						h.Type("button"),
						h.Class("inline-flex items-center p-2 ml-1 text-sm text-gray-500 rounded-lg lg:hidden hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200 dark:text-gray-400 dark:hover:bg-gray-700 dark:focus:ring-gray-600"),
						h.Aria("controls", "mobile-menu-2"),
					),
				),
				h.Div(
					h.Class("hidden justify-between items-center w-full lg:flex lg:w-auto lg:order-1"),
					h.ID("mobile-menu-2"),
					h.Ul(
						h.Class("flex flex-col mt-4 font-medium lg:flex-row lg:space-x-8 lg:mt-0"),
						NavLink("/", "Home", currentPath),
						NavLink("/about", "About", currentPath),
						NavLink("/posts", "Job Posts", currentPath),
					),
				),
			),
		),

		h.Nav(
			c.Classes{"text-xl hidden justify-between": true},
			h.Div(
				h.Class("flex items-center relative"),
				h.Img(
					h.Src("/public/test-logo.svg"),
					h.Height("40px"),
					h.Width("40px"),
				),
				NavbarLink("/", "Home", currentPath),
				NavbarLink("/about", "About", currentPath),
				NavbarLink("/posts", "Job Posts", currentPath),
			),
			h.FormEl(
				c.Classes{"select-none": true},
				g.Attr("onchange", "utils.setTheme(event)"),
				h.Label(
					c.Classes{"flex items-center": true},
					h.Span(
						h.Class("sr-only"),
						h.Span(
							g.Text("theme"),
						),
					),
					h.Select(
						c.Classes{"ml-2 border-none dark:bg-slate-400": true},
						h.Name("themepicker"),
						g.Group(g.Map(config.ThemeOptions, func(option models.ThemeOption) g.Node {
							return h.Option(
								h.Value(option.Value),
								g.If(option.Selected, h.Selected()),
								g.Text(option.Label),
							)
						})),
					),
				),
			),
		),
	)
}

func NavbarLink(href, name, currentPath string) g.Node {
	return h.A(
		h.Href(href),
		c.Classes{
			"p-2":       true,
			"is-active": currentPath == href,
		},
		g.Text(name),
	)
}

func NavLink(href, name, currentPath string) g.Node {
	return h.Li(
		h.A(
			h.Href(href),
			c.Classes{
				"block py-2 pr-4 pl-3 border-b border-gray-100 hover:bg-gray-50 lg:hover:bg-transparent lg:border-0 lg:hover:text-primary-700 lg:p-0 lg:dark:hover:text-white dark:hover:bg-gray-700 dark:hover:text-white lg:dark:hover:bg-transparent dark:border-gray-700": true,
				"text-gray-700 dark:text-gray-400": currentPath != href,
				"text-white dark:text-gray-200":    currentPath == href,
			},
			h.Aria("current", "page"),
			g.Text(name),
		),
	)
}
