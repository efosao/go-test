package controllers

import (
	_ "embed"
	"net/url"
	"strings"
	"vauntly/components"
	"vauntly/models"
	utils "vauntly/utils"

	"github.com/labstack/echo/v4"
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
	"gorm.io/gorm/clause"
)

func Home(c echo.Context) error {
	config, error := utils.GetConfig(c)
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
	Tags := <-tagsChan

	selectedTagMap := map[string]bool{}
	for _, selectedTag := range selectedTags {
		selectedTagMap[selectedTag] = true
	}

	updatedTags := []models.Tag{}
	for _, tag := range Tags {
		if selectedTagMap[tag.Name] {
			tag.Selected = true
		}
		updatedTags = append(updatedTags, tag)
	}

	return PostsPage(config, posts, updatedTags, selectedTagsString, 0, "Search the best TECH Jobs today >>").Render(c.Response().Writer)
}

func HomePage(config *models.Config) g.Node {
	return components.Layout("Home", config,
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
