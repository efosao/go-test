package controllers

import (
	"encoding/json"
	"fmt"
	"vauntly/components"
	"vauntly/models"
	utils "vauntly/utils"

	"github.com/labstack/echo/v4"
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

func GetAbout(c echo.Context) error {
	config, error := utils.GetConfig(c)
	if error != nil {
		return error
	}
	return AboutPage(config).Render(c.Response().Writer)
}

func PostAbout(c echo.Context) error {
	var body struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	body.Email = c.FormValue("email")
	body.Name = c.FormValue("name")
	println("Email: ", body.Email)
	println("Name: ", body.Name)

	return GetAbout(c)
}

func AboutPage(config *models.Config) g.Node {
	LoadTags()

	type Option struct {
		Label string `json:"label"`
		Value string `json:"value"`
	}

	options := make([]Option, len(Tags))

	for index, element := range Tags {
		options[index] = Option{
			Label: element.Name,
			Value: element.Name,
		}
	}

	tagStr := ""
	if tagsBytes, err := json.Marshal(options); err == nil {
		tagStr = string(tagsBytes)
	} else {
		fmt.Println(err)
	}

	return components.Layout("About 2.2", config,
		h.Section(
			h.Class("my-4"),
			h.Div(
				h.Class("mx-auto max-w-screen-xl"),
				h.H3(
					h.Class("text-3xl leading-9 font-extrabold tracking-tight text-gray-900 sm:text-4xl sm:leading-10 pointer-events-none"),
					g.Text("Lit Web-Components === ❤️❤️❤️"),
				),
				// g.Raw("<x-greeting count=5></x-greeting>"),
				// g.Raw("<x-greeting count=15></x-greeting>"),
				h.Div(
					h.Class("my-4 flex flex-col gap-2"),
					g.Raw("<test-rc></test-rc>"),
					g.Raw("<test-rc></test-rc>"),
					g.Raw(fmt.Sprintf("<react-select options='%s'></react-select>", tagStr)),
				),
				h.Div(
					h.ID("post-list"),
					h.Class("mt-4 border rounded-md p-4"),
				),
			),
			getAlpineTest(),
			getAlpineTest(),
			getAlpineTest(),
			getAlpineTest(),
			h.Button(
				c.Classes{"button": true},
				h.ID("dialog_button"),
				g.Text("Show test dialog"),
			),
			h.Dialog(
				h.ID("dialog"),
				c.Classes{"rounded-xl opacity-50 shadow-xl shadow-slate-800 fade-in-bottom overflow-hidden": true},
				h.Div(
					c.Classes{"rounded-xl p-4 min-w-80 w-96 min-h-52 text-black border-2 border-slate-300 relative": true},
					h.Button(
						c.Classes{"close absolute right-2 top-1 bg-[transparent!important]": true},
						h.Img(h.Src("/public/images/close.svg"), h.Height("32"), h.Width("32")),
					),
					h.FormEl(
						c.Classes{"flex flex-col gap-2": true},
						h.Method("post"),
						h.Div(
							c.Classes{"mt-8": true},
							g.Text("This modal dialog has a groovy backdrop!"),
						),
						h.Input(
							h.Type("text"), h.Name("name"), g.Attr("Autofocus"), g.Attr("Autocapitalize"), g.Attr("Autocomplete", "name"), h.Placeholder("Enter your name"), h.Required()),
						h.Input(h.Type("email"), h.Name("email"), g.Attr("Autocomplete", "email"), h.Placeholder("Enter your email"), h.Required()),
						h.Div(
							c.Classes{"flex justify-end gap-2": true},
							h.Button(
								c.Classes{"button bg-transparent close": true},
								h.Type("button"),
								h.Value("cancel"),
								g.Attr("formnovalidate", ""),
								g.Attr("formmethod", "dialog"),
								g.Text("Cancel"),
							),
							h.Button(
								c.Classes{"button btn-apply": true},
								h.Type("submit"),
								g.Text("Submit"),
							),
						),
					),
				),
			),
		),
	)
}

func getAlpineTest() g.Node {
	return h.Div(
		h.Class("flex flex-col gap-2 comp__alpine_test"),
		// g.Attr("x-data", "{ open: false }"),
		g.Attr("up-data", "{ \"open\": false }"),
		h.Label(
			h.Class("text-gray-700 text-xl font-extrabold"),
			h.DataAttr("text", "$foo"),
		),
		h.Button(
			h.Class("text-white p-2 rounded-md transition-colors bg-red-500"),
			g.Attr("@click", "open = !open"),
			g.Text("Toggle block w/h Alpine"),
			// g.Attr("x-bind:class", "open ? 'bg-red-500' : 'bg-black'"),
		),
		h.Div(
			h.Class("bg-pink-200 input p-8 border border-gray-300 dark:bg-pink-600 rounded-md comp__view hidden"),
			// g.Attr("x-show", "open"),
			// g.Attr("x-transition"),
			// g.Attr("x-bind:class", "open ? 'block' : 'hidden'"),
			h.P(
				h.Class("text-black"),
				g.Text("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi porta eu ligula eu sagittis. Donec porta turpis eget euismod congue. Integer nisi arcu, mattis vitae elit quis, dapibus eleifend nulla. Suspendisse posuere ultrices dictum. Morbi in quam interdum, maximus justo at, sollicitudin odio. Etiam ut tempor leo, eu ultricies ipsum. Sed metus lacus, pharetra in ultricies ac, mollis non nisl. Quisque tempus urna et massa facilisis dapibus nec non ex. Nam felis orci, varius quis pulvinar sit amet, placerat et dolor. Duis sed lectus pulvinar, placerat est sit amet, scelerisque tellus. Sed fringilla consectetur dui congue pulvinar. Mauris tempus pharetra enim eu congue. Cras ipsum magna, porta a tellus a, malesuada eleifend neque. Maecenas dictum felis eleifend dignissim tristique. Aliquam rhoncus ultrices est, eu mattis ligula."),
			),
			h.P(
				h.Class("text-black"),
				g.Text("Sed et fringilla diam. Vivamus hendrerit sodales nunc et dictum. Nulla laoreet consequat tellus, elementum ultricies tellus dignissim ac. Donec dui lectus, vehicula eget porttitor eu, volutpat id augue. Nullam odio ipsum, placerat a felis quis, hendrerit dictum orci. Donec egestas justo vel magna semper lobortis. In dictum blandit metus quis dignissim. Cras tristique gravida ipsum sed lacinia."),
			),
		),
	)
}
