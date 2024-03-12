package controllers

import (
	_ "embed"
	"vauntly/components"
	utils "vauntly/utils"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	config, error := utils.GetConfig(c)
	if error != nil {
		return error
	}

	return components.Login(config).Render(c.Response().Writer)
}
