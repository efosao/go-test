package controllers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"vauntly/components"
	"vauntly/models"
	utils "vauntly/utils"

	"github.com/labstack/echo/v4"
	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
	"gorm.io/gorm/clause"
)

func GetPosts(c echo.Context) error {
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

	return PostsPage(config, posts, updatedTags, selectedTagsString, 0, "Posts").Render(c.Response().Writer)
}

func PostsPage(config *models.Config, posts []models.Post, tags []models.Tag, selectedTags string, page int, title string) g.Node {
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

	return components.Layout(title, config,
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
					components.PostSearchResults(posts, selectedTags, nextPage),
				),
			),
		),
	)
}
