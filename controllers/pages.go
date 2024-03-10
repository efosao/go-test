package controllers

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"vauntly/models"

	"github.com/labstack/echo/v4"
	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
	"gorm.io/gorm/clause"
)

var tags = []models.Tag{}
var cacheHash = ""

func SetupCacheHash(hashValue string) {
	if len(hashValue) < 8 {
		cacheHash = hashValue
		return
	}
	cacheHash = hashValue[0:8]
}

func LoadTags() []models.Tag {
	if len(tags) == 0 {
		models.DB.Raw(`
			SELECT unnest(tags) AS name, count(*)::text AS count
			FROM posts
			WHERE published_at IS NOT NULL
			GROUP by name
			ORDER BY count(*) DESC;
		`).Scan(&tags)
	}
	return tags
}

func getConfig(c echo.Context) (*Config, error) {
	themeCookie, _ := c.Cookie("theme")
	theme := "system"
	if themeCookie != nil {
		theme = themeCookie.Value
	}

	themeOptions := []models.ThemeOption{
		{Value: "light", Label: "🌞", Selected: theme == "light"},
		{Value: "dark", Label: "🌘", Selected: theme == "dark"},
		{Value: "system", Label: "🌎", Selected: theme != "light" && theme != "dark"},
	}

	config := &Config{
		path:         c.Request().URL.Path,
		theme:        theme,
		themeOptions: themeOptions,
	}
	return config, nil
}

func GetHome(c echo.Context) error {
	config, error := getConfig(c)
	if error != nil {
		return error
	}
	return HomePage(config).Render(c.Response().Writer)
}

func HomePage(config *Config) g.Node {
	return Layout("Home", config,
		h.Section(
			c.Classes{"my-4": true},
			h.Div(
				c.Classes{"mx-auto max-w-screen-xl": true},
				h.H3(
					h.Class("text-3xl leading-9 font-extrabold tracking-tight text-gray-900 sm:text-4xl sm:leading-10 pointer-events-none"),
					g.Text("Welcome to the job board"),
				),
				h.P(
					c.Classes{"mt-4 text-lg leading-7 text-gray-500": true},
					g.Text("This is a job board for the modern web."),
				),
			),
		),
	)
}

func GetAbout(c echo.Context) error {
	config, error := getConfig(c)
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

func AboutPage(config *Config) g.Node {
	LoadTags()

	type Option struct {
		Label string `json:"label"`
		Value string `json:"value"`
	}

	options := make([]Option, len(tags))

	for index, element := range tags {
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

	return Layout("About 2.2", config,
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

func GetPosts(c echo.Context) error {
	config, error := getConfig(c)
	if error != nil {
		return error
	}
	selectedTagsString := c.QueryParam("tags")
	var selectedTags []string
	if selectedTagsString != "" {
		selectedTags = strings.Split(selectedTagsString, ",")
	}
	unescapedSelectedTags := []string{}
	for _, selectedTag := range selectedTags {
		escapedTag, err := url.QueryUnescape(selectedTag)
		if err == nil {
			// This is a hack to fix the fact that the "c++" tag is not being unescaped properly
			if escapedTag == "c  " {
				escapedTag = strings.ReplaceAll(escapedTag, " ", "+")
			}
			unescapedSelectedTags = append(unescapedSelectedTags, escapedTag)
		}
	}

	postsChan := make(chan []models.Post)
	tagsChan := make(chan []models.Tag)

	go (func(p chan []models.Post) {
		posts := []models.Post{}

		if len(unescapedSelectedTags) > 0 {
			queryInputTags := "{" + strings.Join(unescapedSelectedTags, ",") + "}"
			models.DB.Select("ID", "CompanyName", "Location", "Tags", "Thumbnail", "Title", "PublishedAt", "CreatedAt").Where("tags @> ?", queryInputTags).Where("published_at IS NOT NULL").Order(clause.OrderByColumn{Column: clause.Column{Name: "published_at"}, Desc: true}).Limit(10).Find(&posts)
			p <- posts
			return
		} else {
			models.DB.Select("ID", "CompanyName", "Location", "Tags", "Thumbnail", "Title", "PublishedAt", "CreatedAt").Where("published_at IS NOT NULL").Order(clause.OrderByColumn{Column: clause.Column{Name: "published_at"}, Desc: true}).Limit(10).Find(&posts)
			p <- posts
		}
	})(postsChan)

	go (func(t chan []models.Tag) {
		t <- LoadTags()
	})(tagsChan)

	posts := <-postsChan
	tags := <-tagsChan

	selectedTagMap := map[string]bool{}
	for _, selectedTag := range selectedTags {
		selectedTagMap[selectedTag] = true
	}

	updatedTags := []models.Tag{}
	for _, tag := range tags {
		if selectedTagMap[tag.Name] {
			tag.Selected = true
		}
		updatedTags = append(updatedTags, tag)
	}

	return PostsPage(config, posts, updatedTags, selectedTagsString, 0).Render(c.Response().Writer)
}

func PostsPage(config *Config, posts []models.Post, tags []models.Tag, selectedTags string, page int) g.Node {
	nextPage := page + 1

	type Option struct {
		Label    string `json:"label"`
		Value    string `json:"value"`
		Selected bool   `json:"selected"`
	}

	options := make([]Option, len(tags))

	for index, element := range tags {
		options[index] = Option{
			Label:    element.Name,
			Value:    element.Name,
			Selected: element.Selected,
		}
	}

	tagStr := ""
	if tagsBytes, err := json.Marshal(options); err == nil {
		tagStr = string(tagsBytes)
	} else {
		fmt.Println(err)
	}

	return Layout("Posts", config,
		h.Section(
			hx.History("false"), // disable htmx caching for this page
			h.Div(
				h.Class("h-9"),
				g.Raw(fmt.Sprintf("<react-select options='%s'></react-select>", tagStr)),
			),
			h.Div(
				h.ID("post-list"),
				h.Class("mt-4"),
				h.Div(
					Posts(posts, selectedTags, nextPage),
				),
			),
		),
	)
}

func Posts(posts []models.Post, selectedTags string, nextPage int) g.Node {
	nextPgUrl := fmt.Sprint("/partials/posts/search/", nextPage)
	if selectedTags != "" {
		nextPgUrl = fmt.Sprint("/partials/posts/search/", nextPage, "?tags=", selectedTags)
	}

	return h.Div(
		g.Group(g.Map(posts, func(post models.Post) g.Node {
			return Post(post)
		})),
		h.Div(
			g.Attr("hx-post", nextPgUrl),
			g.Attr("hx-swap", "outerHTML"),
			g.Attr("hx-trigger", "revealed"),
			h.ID("nextPageLoaderId_"+fmt.Sprint(nextPage)),
			h.Class("htmx-indicator flex flex-col gap-2 items-center justify-center"),
			h.Div(
				h.Class("bg-orange-200 dark:bg-slate-800 rounded-md h-36 w-full"),
			),
			h.Div(
				h.Class("bg-orange-200 dark:bg-slate-800 rounded-md h-36 w-full"),
			),
			h.Div(
				h.Class("bg-orange-200 dark:bg-slate-800 rounded-md h-36 w-full"),
			),
			h.Div(
				h.Class("bg-orange-200 dark:bg-slate-800 rounded-md h-36 w-full"),
			),
			h.Div(
				h.Class("bg-orange-200 dark:bg-slate-800 rounded-md h-36 w-full"),
			),
			h.Div(
				h.Class("bg-orange-200 dark:bg-slate-800 rounded-md h-36 w-full"),
			),
		),
	)
}

func Post(post models.Post) g.Node {
	class := "search_row w-full group relative mb-2 rounded-md border-0 border-pink-200 bg-pink-100 dark:border-prussian-blue-900 dark:bg-slate-700 dark:text-white"
	if post.IsPinned() {
		class = "search_row group relative mb-2 rounded-md border-0 border-orange-200 bg-orange-200 text-black dark:border-slate-700 dark:bg-pink-700 dark:text-white"
	}

	return h.Div(
		h.Class(class),
		g.Attr("onclick", "utils.toggleOpenState('cbx"+post.ID+"', 'desc"+post.ID+"')"),
		h.Div(
			c.Classes{"cursor-pointer flex h-32 items-center space-x-2 px-2": true},
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
							c.Classes{"inline cursor-pointer rounded-md bg-white px-2 mr-1 font-semibold text-pink-950 transition-colors duration-300 hover:bg-blue-100 hover:text-black my-[2px]": true},
							g.Text(tag),
						)
					})),
				),
			),
			h.Span(
				c.Classes{"m-2": true},
				g.Text(post.TimeSinceCreated()),
			),
			h.Span(
				g.Attr("onclick", "utils.halt(event)"),
				h.Button(
					c.Classes{"btn-apply done": true},
					g.Text("Applied"),
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

type Config struct {
	path         string
	theme        string
	themeOptions []models.ThemeOption
}

func GetPostDetail(c echo.Context) error {
	id := c.Param("id")
	post := &models.Post{}
	if err := models.DB.Select("ID", "Title", "Description").Where(&models.Post{ID: id}).First(&post).Error; err != nil {
		return err
	}

	return PostDetailPage(post).Render(c.Response().Writer)
}

func PostDetailPage(post *models.Post) g.Node {
	return h.Section(
		h.Class("my-4"),
		h.Div(
			h.Class("flex flex-col items-center gap-2 max-w-screen-xl"),
			h.H3(
				c.Classes{"text-3xl leading-9 font-extrabold tracking-tight text-gray-900 sm:text-4xl sm:leading-10": true},
				g.Text(post.Title),
			),
			h.Article(
				c.Classes{"mt-4 text-lg leading-7 text-gray-500": true},
				g.Raw(post.GetDescription()),
			),
		),
	)
}

func Layout(title string, config *Config, children g.Node) g.Node {
	releaseHash := fmt.Sprintf("?v=%s", cacheHash)
	return h.Doctype(
		h.HTML(
			h.Class("smooth-scroll "+config.theme),
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
				h.Class("dark:bg-slate-400"),
				Navbar(config),
				h.Div(
					h.Class("bg-slate-600 h-14 p-2"),
					h.Div(
						h.Class("mx-auto max-w-7xl"),
						g.Raw("<app-bar><div class='h-10 bg-black rounded-md'></div></app-bar>"),
					),
				),
				h.H1(
					h.ID("page-title"),
					h.Class("container mx-auto my-4 max-w-5xl text-3xl font-extrabold text-black dark:text-black pointer-events-none"),
					g.Text(title),
				),
				h.Div(
					h.Class("container mx-auto max-w-5xl"),
					children,
				),
				h.Div(
					h.Class("flex justify-center gap-2 mt-4 mb-2"),
					h.Span(
						h.Class("text-md font-bold text-gray-900 bg-slate-300 p-2 rounded-md"),
						g.Text("Release"),
						h.Span(
							h.Class("text-md font-bold text-gray-700 p-1 ml-2 bg-slate-400 rounded-md"),
							g.Text(cacheHash),
						),
					),
				),
			),
		),
	)
}

func Navbar(config *Config) g.Node {
	currentPath := config.path
	return h.Div(
		g.Raw(`<nav class="bg-white border-gray-200 px-4 lg:px-6 py-2.5 dark:bg-gray-800">
        <div class="flex flex-wrap justify-between items-center mx-auto max-w-screen-xl">
            <a href="/" class="flex items-center">
                <img src="/public/test-logo.svg" class="mr-3 h-6 sm:h-9" alt="Flowbite Logo" />
                <span class="self-center text-xl font-semibold whitespace-nowrap dark:text-white">Vauntly</span>
            </a>
            <div class="flex items-center lg:order-2">
                <a href="#" class="text-gray-800 dark:text-white hover:bg-gray-50 focus:ring-4 focus:ring-gray-300 font-medium rounded-lg text-sm px-4 lg:px-5 py-2 lg:py-2.5 mr-2 dark:hover:bg-gray-700 focus:outline-none dark:focus:ring-gray-800">Log in</a>
                <a href="#" class="text-white bg-primary-700 hover:bg-primary-800 focus:ring-4 focus:ring-primary-300 font-medium rounded-lg text-sm px-4 lg:px-5 py-2 lg:py-2.5 mr-2 dark:bg-primary-600 dark:hover:bg-primary-700 focus:outline-none dark:focus:ring-primary-800">Get started</a>
                <button data-collapse-toggle="mobile-menu-2" type="button" class="inline-flex items-center p-2 ml-1 text-sm text-gray-500 rounded-lg lg:hidden hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200 dark:text-gray-400 dark:hover:bg-gray-700 dark:focus:ring-gray-600" aria-controls="mobile-menu-2" aria-expanded="false">
                    <span class="sr-only">Open main menu</span>
                    <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M3 5a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 15a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd"></path></svg>
                    <svg class="hidden w-6 h-6" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"></path></svg>
                </button>
            </div>
            <div class="hidden justify-between items-center w-full lg:flex lg:w-auto lg:order-1" id="mobile-menu-2">
                <ul class="flex flex-col mt-4 font-medium lg:flex-row lg:space-x-8 lg:mt-0">
                    <li>
                        <a href="/" class="block py-2 pr-4 pl-3 text-white rounded bg-primary-700 lg:bg-transparent lg:text-primary-700 lg:p-0 dark:text-white" aria-current="page">Home</a>
                    </li>
                    <li>
                        <a href="/about" class="block py-2 pr-4 pl-3 text-gray-700 border-b border-gray-100 hover:bg-gray-50 lg:hover:bg-transparent lg:border-0 lg:hover:text-primary-700 lg:p-0 dark:text-gray-400 lg:dark:hover:text-white dark:hover:bg-gray-700 dark:hover:text-white lg:dark:hover:bg-transparent dark:border-gray-700">About</a>
                    </li>
                    <li>
                        <a href="/posts" class="block py-2 pr-4 pl-3 text-gray-700 border-b border-gray-100 hover:bg-gray-50 lg:hover:bg-transparent lg:border-0 lg:hover:text-primary-700 lg:p-0 dark:text-gray-400 lg:dark:hover:text-white dark:hover:bg-gray-700 dark:hover:text-white lg:dark:hover:bg-transparent dark:border-gray-700">Posts</a>
                    </li>
                </ul>
            </div>
        </div>
    </nav>`),
		h.Nav(
			hx.Boost("true"),
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
						g.Group(g.Map(config.themeOptions, func(option models.ThemeOption) g.Node {
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
