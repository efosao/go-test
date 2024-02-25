package controllers

import (
	"fmt"
	"gofiber/models"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"gorm.io/gorm/clause"
)

func PostSearchResultsPage(w http.ResponseWriter, r *http.Request) {
	pageStr := r.PathValue("page")
	page := 0
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	nextPage := page + 1
	offset := (nextPage - 1) * 10

	selectedTagsStr := ""
	if r.Method == "GET" {
		selectedTagsStr = r.URL.Query().Get("tags")
	} else {
		r.ParseForm()
		tags := r.Form["tags"]
		selectedTagsStr = strings.Join(tags, ",")
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

	w.Header().Set("Content-Type", "text/html")

	if selectedTagsStr == "" {
		if page == 0 {
			w.Header().Set("HX-Push-Url", "/posts")
		}
		models.DB.Select("ID", "CompanyName", "Location", "Tags", "Thumbnail", "Title", "PublishedAt", "CreatedAt").Where("published_at IS NOT NULL").Order(clause.OrderByColumn{Column: clause.Column{Name: "published_at"}, Desc: true}).Offset(offset).Limit(10).Find(&posts)

		if len(posts) == 0 {
			fmt.Fprintln(w, "")
			return
		}

		Posts(posts, selectedTagsStr, nextPage).Render(w)
	} else {
		if page == 0 {
			w.Header().Set("HX-Push-Url", fmt.Sprintf("/posts?tags=%s", selectedTagsStr))
		}
		models.DB.Select("ID", "CompanyName", "Location", "Tags", "Thumbnail", "Title", "PublishedAt", "CreatedAt").Where("tags @> ?", queryInputTags).Where("published_at IS NOT NULL").Order(clause.OrderByColumn{Column: clause.Column{Name: "published_at"}, Desc: true}).Offset(offset).Limit(10).Find(&posts)

		if len(posts) == 0 {
			fmt.Fprintln(w, "")
			return
		}

		Posts(posts, selectedTagsStr, nextPage).Render(w)
	}
}
