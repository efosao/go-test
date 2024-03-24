package controllers

import (
	"vauntly/components"
	utils "vauntly/utils"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	config, error := utils.GetConfig(c)
	config.ShowNav = false
	if error != nil {
		return error
	}

	return components.Login(config).Render(c.Response().Writer)
}
