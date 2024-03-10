package controllers

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"vauntly/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm/clause"
)

func PostSearchResultsPage(c echo.Context) error {
	pageStr := c.Param("page")
	page := 0
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	nextPage := page + 1
	offset := (nextPage - 1) * 10

	selectedTagsStr := ""
	if c.Request().Method == "GET" {
		selectedTagsStr = c.QueryParams().Get("tags")
	} else {
		// TODO: learn this Go pattern
		if params, err := c.FormParams(); err != nil {
			fmt.Println("Error parsing form params", err)
		} else {
			tags := params["tags"]
			selectedTagsStr = strings.Join(tags, ",")
		}
	}

	posts := []models.Post{}
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

	w := c.Response().Writer

	w.Header().Set("Content-Type", "text/html")

	host := c.Request().Header.Get("Referer")
	if strings.Contains(host, "?") {
		host = host[:strings.Index(host, "?")]
	}

	if selectedTagsStr == "" {
		if page == 0 {
			w.Header().Set("HX-Push-Url", host)
		}
		models.DB.Select("ID", "CompanyName", "Location", "Tags", "Thumbnail", "Title", "PublishedAt", "CreatedAt").Where("published_at IS NOT NULL").Order(clause.OrderByColumn{Column: clause.Column{Name: "published_at"}, Desc: true}).Offset(offset).Limit(10).Find(&posts)

		if len(posts) == 0 {
			fmt.Fprintln(w, "")
			return nil
		}

		Posts(posts, selectedTagsStr, nextPage).Render(w)
	} else {
		if page == 0 {
			w.Header().Set("HX-Push-Url", fmt.Sprintf(host+"?tags=%s", selectedTagsStr))
		}
		models.DB.Select("ID", "CompanyName", "Location", "Tags", "Thumbnail", "Title", "PublishedAt", "CreatedAt").Where("tags @> ?", queryInputTags).Where("published_at IS NOT NULL").Order(clause.OrderByColumn{Column: clause.Column{Name: "published_at"}, Desc: true}).Offset(offset).Limit(10).Find(&posts)

		if len(posts) == 0 {
			fmt.Fprintln(w, "")
			return nil
		}

		return Posts(posts, selectedTagsStr, nextPage).Render(c.Response().Writer)
	}

	return nil
}
