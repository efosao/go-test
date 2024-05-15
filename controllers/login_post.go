package controllers

import (
	"fmt"
	"vauntly/components"
	"vauntly/models"
	utils "vauntly/utils"

	"github.com/labstack/echo/v4"
)

func LoginPost(c echo.Context) error {
	config, error := utils.GetConfig(c)
	if error != nil {
		return error
	}
	if params, err := c.FormParams(); err != nil {
		fmt.Println("Error parsing form params", err)
	} else {
		email := params["email"]
		password := params["password"]
		fmt.Println("Email: ", email)
		fmt.Println("Password: ", password)

		if email[0] == "efosao@gmail.com" {
			fmt.Println("Email is not valid")
			res := c.Response()
			res.WriteHeader(422)
			return components.Login(config, models.LoginProps{
				Email:    email[0],
				Password: password[0],
				EmailErr: "Email is not valid",
			}).Render(res.Writer)
		}

		// if isValidationRequest, return the form without submitting
		isValidationRequest := len(c.Request().Header.Get("X-Up-Validate")) > 0
		if isValidationRequest {
			return components.Login(config, models.LoginProps{
				Email:    email[0],
				Password: password[0],
			}).Render(c.Response().Writer)
		}
	}

	res := c.Response()
	// close the login layer
	res.Header().Set("X-Up-Dismiss-Layer", "true")
	res.Write([]byte("You are logged in"))
	return nil
}
