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
					h.Class("flex justify-between items-center w-auto"),
					h.ID("mobile-menu-2"),
					h.Ul(
						h.Class("flex font-medium flex-row gap-4"),
						NavLink("/", "Home", currentPath),
						NavLink("/about", "About", currentPath),
						NavLink("/posts", "Job Posts", currentPath),
					),
				),
				h.Div(
					h.Class("hidden lg:flex items-center"),
					h.A(
						h.Href("/login"),
						h.Class("text-gray-800 dark:text-white hover:bg-gray-50 focus:ring-4 focus:ring-gray-300 font-medium rounded-lg text-sm px-4 lg:px-5 py-2 lg:py-2.5 mr-2 dark:hover:bg-gray-700 focus:outline-none dark:focus:ring-gray-800"),
						g.Text("Log in"),
					),
					h.A(
						h.Href("#"),
						h.Class("hidden text-white bg-primary-700 hover:bg-primary-800 focus:ring-4 focus:ring-primary-300 font-medium rounded-lg text-sm px-4 lg:px-5 py-2 lg:py-2.5 mr-2 dark:bg-primary-600 dark:hover:bg-primary-700 focus:outline-none dark:focus:ring-primary-800"),
						g.Text("Get started"),
					),
					h.Button(
						h.DataAttr("collapse-toggle", "mobile-menu-2"),
						h.Type("button"),
						h.Class("inline-flex items-center p-2 ml-1 text-sm text-gray-500 rounded-lg lg:hidden hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200 dark:text-gray-400 dark:hover:bg-gray-700 dark:focus:ring-gray-600"),
						h.Aria("controls", "mobile-menu-2"),
					),
				),
				h.Button(
					h.Class("p-4 relative hover:bg-transparent dark:hover:bg-transparent"),
					h.Type("button"),
					g.Attr("onclick", "utils.toggleCurrentTheme()"),
					h.Span(
						h.ID("dark-mode-toggle"),
						c.Classes{
							"absolute top-1 left-1": true,
							"hidden":                config.Theme == "dark",
						},
						// sun icon
						g.Raw(`<svg class="w-6 h-6 text-slate-900" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
								<path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5V3m0 18v-2M7.05 7.05 5.636 5.636m12.728 12.728L16.95 16.95M5 12H3m18 0h-2M7.05 16.95l-1.414 1.414M18.364 5.636 16.95 7.05M16 12a4 4 0 1 1-8 0 4 4 0 0 1 8 0Z"/>
							</svg>
							`),
					),
					h.Span(
						h.ID("light-mode-toggle"),
						c.Classes{
							"absolute top-1 left-1": true,
							"hidden":                config.Theme != "dark",
						},
						// moon icon
						g.Raw(`<svg class="w-6 h-6 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
							<path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 21a9 9 0 0 1-.5-17.986V3c-.354.966-.5 1.911-.5 3a9 9 0 0 0 9 9c.239 0 .254.018.488 0A9.004 9.004 0 0 1 12 21Z"/>
						</svg>				  
						`),
					),
				),
			),
		),
	)
}

func NavLink(href, name, currentPath string) g.Node {
	return h.Li(
		h.A(
			h.Href(href),
			c.Classes{
				"block p-2 border-gray-100 hover:bg-gray-50 lg:hover:bg-transparent lg:border-0 lg:hover:text-primary-700 lg:p-0 lg:dark:hover:text-white dark:hover:bg-gray-700 dark:hover:text-white lg:dark:hover:bg-transparent dark:border-gray-700": true,
				"text-gray-500 dark:text-gray-400": currentPath != href,
				"text-gray-900 dark:text-gray-200": currentPath == href,
			},
			h.Aria("current", "page"),
			g.Text(name),
		),
	)
}
