package controllers

import (
	"fmt"
	"gofiber/models"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
	"gorm.io/gorm/clause"
)

func GetHome(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title":        "Hello, World!",
		"Description":  "Find the latest job posts in the tech industry.",
		"ThemeOptions": c.Locals("ThemeOptions"),
	}, "layouts/main")
}

func GetPosts(c *fiber.Ctx) error {
	selectedTagsStr := c.Query("tags")
	selectedTags := strings.Split(selectedTagsStr, ",")
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
		if len(selectedTags) > 0 {
			queryInputTags := "{" + strings.Join(unescapedSelectedTags, ",") + "}"
			models.DBConn.Select("ID", "CompanyName", "Location", "Tags", "Thumbnail", "Title", "PublishedAt", "CreatedAt").Where("tags @> ?", queryInputTags).Where("published_at IS NOT NULL").Order(clause.OrderByColumn{Column: clause.Column{Name: "published_at"}, Desc: true}).Limit(10).Find(&posts)
			p <- posts
			return
		} else {
			models.DBConn.Select("ID", "CompanyName", "Location", "Tags", "Thumbnail", "Title", "PublishedAt", "CreatedAt").Where("published_at IS NOT NULL").Order(clause.OrderByColumn{Column: clause.Column{Name: "published_at"}, Desc: true}).Limit(10).Find(&posts)
			p <- posts
		}
	})(postsChan)

	go (func(t chan []models.Tag) {
		tags := []models.Tag{}
		models.DBConn.Raw(`
			SELECT unnest(tags) AS name, count(*)::text AS count
			FROM posts
			WHERE published_at IS NOT NULL
			GROUP by name
			ORDER BY count(*) DESC;
		`).Scan(&tags)
		t <- tags
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

	cookie := new(models.Cookie)
	if err := c.CookieParser(cookie); err != nil {
		return err
	}

	data := fiber.Map{
		"Description":     "Find the latest job posts in the tech industry.",
		"Page":            "1",
		"Posts":           posts,
		"SelectedTagsStr": selectedTagsStr,
		"Theme":           cookie.Theme,
		"Tags":            updatedTags,
		"ThemeOptions":    c.Locals("ThemeOptions"),
		"Title":           "Job Posts",
	}

	return c.Render("posts", data, "layouts/main")
}

type Config struct {
	path         string
	theme        string
	themeOptions []models.ThemeOption
}

func GetAbout(w http.ResponseWriter, r *http.Request) {
	if themeOptions, ok := r.Context().Value("themeOption").([]models.ThemeOption); ok {
		config := &Config{
			path:         r.URL.Path,
			theme:        r.Context().Value("theme").(string),
			themeOptions: themeOptions,
		}
		fmt.Println("config", config)
		AboutPage(config).Render(w)
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func AboutPage(config *Config) g.Node {
	return Layout("About", config,
		h.Section(
			c.Classes{"my-4": true},
			h.Div(
				c.Classes{"mx-auto max-w-screen-xl": true},
				h.H3(
					c.Classes{"text-3xl leading-9 font-extrabold tracking-tight text-gray-900 sm:text-4xl sm:leading-10": true},
					g.Text("Lit Web-Components === ❤️"),
				),
				g.Raw("<x-greeting count=5></x-greeting>"),
				g.Raw("<x-greeting count=15></x-greeting>"),
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
					c.Classes{"rounded-xl p-4 min-w-80 min-h-52 text-black border-2 border-slate-300 relative": true},
					h.Button(
						c.Classes{"close absolute right-2 top-1 bg-[transparent!important]": true},
						g.Attr("onclick", "utils.closeDialog()"),
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

func Layout(title string, config *Config, children g.Node) g.Node {
	return h.Doctype(
		h.HTML(
			c.Classes{config.theme: true},
			hx.Boost("true"),
			h.Lang("en"),
			h.Head(
				h.TitleEl(g.Text(title)),
				h.StyleEl(h.Type("text/css"), g.Raw(".is-active{ font-weight: bold }")),
				h.Link(h.Rel("stylesheet"), h.Href("/public/dist/stylesheet.css")),
				h.Script(h.Src("/public/dist/index.js"), h.Defer()),
			),
			h.Body(
				c.Classes{"max-w-4xl mx-auto dark:bg-slate-500": true},
				Navbar(config),
				h.H1(
					c.Classes{"text-3xl font-extrabold mb-4 text-black dark:text-black": true},
					g.Text(title),
				),
			),
			children,
		),
	)
}

func Navbar(config *Config) g.Node {
	currentPath := config.path
	return h.Nav(
		c.Classes{"m-2 text-xl flex justify-between": true},
		h.Div(
			c.Classes{"flex gap-2": true},
			NavbarLink("/", "Home", currentPath),
			NavbarLink("/posts/", "Job Posts", currentPath),
			NavbarLink("/about/", "About", currentPath),
			g.Text("👀"),
		),
		h.FormEl(
			c.Classes{"select-none": true},
			g.Attr("onchange", "utils.setTheme(event)"),
			h.Label(
				c.Classes{"flex items-center": true},
				h.Span(c.Classes{"text-slate-900 dark:text-slate-900": true}, g.Text("theme")),
				h.Select(
					c.Classes{"ml-2 border-none dark:bg-slate-500": true},
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
	)
}

func NavbarLink(href, name, currentPath string) g.Node {
	return h.A(h.Href(href), c.Classes{"is-active": currentPath == href}, g.Text(name))
}
