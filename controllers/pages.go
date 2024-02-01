package controllers

import (
	"gofiber/models"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
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
