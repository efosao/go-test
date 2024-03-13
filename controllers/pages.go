package controllers

import (
	"vauntly/models"
)

var Tags = []models.Tag{}

func LoadTags() []models.Tag {
	if len(Tags) == 0 {
		models.DB.Raw(`
			SELECT unnest(tags) AS name, count(*)::text AS count
			FROM posts
			WHERE published_at IS NOT NULL
			GROUP by name
			ORDER BY count(*) DESC;
		`).Scan(&Tags)
	}
	return Tags
}
