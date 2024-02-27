package controllers

import (
	"fmt"
	"gofiber/models"
	"net/http"
	"net/url"
	"strings"

	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
	"gorm.io/gorm/clause"
)

var tags = []models.Tag{}

func GetHome(w http.ResponseWriter, r *http.Request) {
	if themeOptions, ok := r.Context().Value(models.ThemeOptionsKey).([]models.ThemeOption); ok {
		config := &Config{
			path:         r.URL.Path,
			theme:        r.Context().Value(models.ThemeKey).(string),
			themeOptions: themeOptions,
		}
		HomePage(config).Render(w)
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func HomePage(config *Config) g.Node {
	return Layout("Home", config,
		h.Section(
			c.Classes{"my-4": true},
			h.Div(
				c.Classes{"mx-auto max-w-screen-xl": true},
				h.H3(
					c.Classes{"text-3xl leading-9 font-extrabold tracking-tight text-gray-900 sm:text-4xl sm:leading-10": true},
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

func GetPosts(w http.ResponseWriter, r *http.Request) {
	selectedTagsString := r.URL.Query().Get("tags")
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
		if len(tags) > 0 {
			t <- tags
			return
		}

		models.DB.Raw(`
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

	config := &Config{
		path:         r.URL.Path,
		theme:        r.Context().Value(models.ThemeKey).(string),
		themeOptions: r.Context().Value(models.ThemeOptionsKey).([]models.ThemeOption),
	}

	PostsPage(config, posts, updatedTags, selectedTagsString, 1).Render(w)
}

func PostsPage(config *Config, posts []models.Post, tags []models.Tag, selectedTags string, page int) g.Node {
	nextPage := page + 1

	return Layout("Posts", config,
		h.Section(
			h.Class("my-4"),
			h.Div(
				h.Class("h-10"),
				h.Select(
					c.Classes{"hidden slim-select": true},
					h.ID("tags"),
					h.Name("tags"),
					hx.Post("/partials/posts/search/0"),
					hx.Target("#post-list"),
					hx.Trigger("change"),
					h.Multiple(),
					h.TabIndex("-1"),
					g.Attr("aria-hidden", "true"),
					g.Group(g.Map(tags, func(tag models.Tag) g.Node {
						return h.Option(
							h.Value(tag.Name),
							g.If(tag.Selected, h.Selected()),
							g.Text(tag.Name),
						)
					})),
				),
			),
			h.Div(
				Posts(posts, selectedTags, nextPage),
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
		h.ID("post-list"),
		h.Class("mt-4"),
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
	class := "search_row group relative mb-2 rounded-md border-0 border-pink-200 bg-pink-100 dark:border-prussian-blue-900 dark:bg-black dark:text-white"
	if post.IsPinned() {
		class = "search_row group relative mb-2 rounded-md border-0 border-orange-200 bg-orange-200 text-black dark:border-slate-700 dark:bg-slate-700 dark:text-white"
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
					g.Attr("hx-get", "/posts/details/"+post.ID),
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

func GetAbout(w http.ResponseWriter, r *http.Request) {
	if themeOptions, ok := r.Context().Value(models.ThemeOptionsKey).([]models.ThemeOption); ok {
		config := &Config{
			path:         r.URL.Path,
			theme:        r.Context().Value(models.ThemeKey).(string),
			themeOptions: themeOptions,
		}
		AboutPage(config).Render(w)
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func AboutPage(config *Config) g.Node {
	return Layout("About 2.2", config,
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

func GetPostDetail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	post := &models.Post{}
	if err := models.DB.Select("ID", "Title", "Description").Where(&models.Post{ID: id}).First(&post).Error; err != nil {
		if err.Error() == "record not found" {
			w.WriteHeader(404)
			w.Write([]byte("Post not found"))
		} else {
			w.WriteHeader(500)
			w.Write([]byte("Internal Server Error"))
		}
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)

	PostDetailPage(post).Render(w)
}

func PostDetailPage(post *models.Post) g.Node {
	return h.Section(
		h.Class("my-4"),
		h.Div(
			c.Classes{"mx-auto max-w-screen-xl": true},
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
				h.Meta(h.Name("viewport"), h.Content("width=device-width, initial-scale=1")),
			),
			h.Body(
				c.Classes{"max-w-4xl mx-auto dark:bg-slate-500": true},
				Navbar(config),
				h.H1(
					h.Class("text-3xl font-extrabold mb-4 text-black dark:text-black mx-2"),
					g.Text(title),
				),
			),
			h.Div(
				h.Class("mx-2"),
				children,
			),
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
			NavbarLink("/about/", "About", currentPath),
			NavbarLink("/posts/", "Job Posts", currentPath),
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
