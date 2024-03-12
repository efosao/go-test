package controllers

import (
	_ "embed"
	"vauntly/models"

	"github.com/labstack/echo/v4"
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

func GetPostDetail(c echo.Context) error {
	id := c.Param("id")
	post := &models.Post{}
	if err := models.DB.Select("ID", "Title", "Description").Where(&models.Post{ID: id}).First(&post).Error; err != nil {
		return err
	}

	return PostDetailPage(post).Render(c.Response().Writer)
}

func PostDetailPage(post *models.Post) g.Node {
	return h.Section(
		h.Class("my-4"),
		h.Div(
			h.Class("flex flex-col items-center gap-2 max-w-screen-xl"),
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
