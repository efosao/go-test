package controllers

import (
	"fmt"
	"gofiber/models"
	"net/url"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func PostSearchResultsPage(c *fiber.Ctx) error {
	type Body struct{ Tags []string }
	if pageStr, err := strconv.Atoi(c.Params("page", "0")); err != nil {
		return err
	} else {
		body := Body{}
		c.BodyParser(&body)
		page := int(pageStr)
		nextPage := page + 1
		offset := (nextPage - 1) * 10
		posts := []models.Post{}
		selectedTagsStr := ""
		if len(body.Tags) > 0 {
			selectedTagsStr = strings.Join(body.Tags, ",")
		} else {
			selectedTagsStr = c.Query("tags")
		}
		unescapedSelectedTags := []string{}
		for _, selectedTag := range strings.Split(selectedTagsStr, ",") {
			escapedTag, err := url.QueryUnescape(selectedTag)
			if err == nil {
				// This is a hack to fix the fact that the "c++" tag is not being unescaped properly
				if escapedTag == "c  " {
					escapedTag = strings.ReplaceAll(escapedTag, " ", "+")
				}
				unescapedSelectedTags = append(unescapedSelectedTags, escapedTag)
			}
		}

		queryInputTags := "{" + strings.Join(unescapedSelectedTags, ",") + "}"

		if selectedTagsStr == "" {
			if page == 0 {
				c.Response().Header.Set("HX-Push-Url", "/posts")
			}
			models.DBConn.Select("ID", "CompanyName", "Location", "Tags", "Thumbnail", "Title", "PublishedAt", "CreatedAt").Where("published_at IS NOT NULL").Order(clause.OrderByColumn{Column: clause.Column{Name: "published_at"}, Desc: true}).Offset(offset).Limit(10).Find(&posts)

			if len(posts) == 0 {
				return c.Send([]byte(""))
			}

			return c.Render("post_list", fiber.Map{
				"Posts": posts,
				"Page":  nextPage,
			})
		} else {
			if page == 0 {
				c.Response().Header.Set("HX-Push-Url", fmt.Sprintf("/posts?tags=%s", selectedTagsStr))
			}
			models.DBConn.Select("ID", "CompanyName", "Location", "Tags", "Thumbnail", "Title", "PublishedAt", "CreatedAt").Where("tags @> ?", queryInputTags).Where("published_at IS NOT NULL").Order(clause.OrderByColumn{Column: clause.Column{Name: "published_at"}, Desc: true}).Offset(offset).Limit(10).Find(&posts)

			if len(posts) == 0 {
				return c.Send([]byte(""))
			}

			return c.Render("post_list", fiber.Map{
				"Posts":           posts,
				"Page":            nextPage,
				"SelectedTagsStr": selectedTagsStr,
			})
		}
	}
}
