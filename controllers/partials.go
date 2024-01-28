package controllers

import (
	"gofiber/models"

	"github.com/gofiber/fiber/v2"
)

func GetPostDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	post := &models.Post{}
	models.DBConn.Select("ID", "Title", "Description").Where(&models.Post{ID: id}).First(&post)
	return c.Render("post_details", fiber.Map{
		"Title":       post.Title,
		"ID":          post.ID,
		"Description": post.GetDescription(),
	})
}
