package handlers

import (
	"net/http"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/labstack/echo"
)

func HomePage(c echo.Context) error {
	/*
		return templates.Page(
			"Home page",
			"", "",
			nil,
			func() templ.Component { return templates.LoginForm() },
			nil,
		).Render(c.Request().Context(), c.Response())
	*/
	return templates.LoginPage().Render(c.Request().Context(), c.Response())
}

func LogIn(c echo.Context) error {
	userNameCookie := new(http.Cookie)
	userNameCookie.Name = "username"
	userName := c.FormValue("username")
	userNameCookie.Value = userName
	userNameCookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(userNameCookie)

	emailCookie := new(http.Cookie)
	emailCookie.Name = "email"
	email := c.FormValue("email")
	emailCookie.Value = email
	emailCookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(emailCookie)

	fs.SessionManager.AddSession(userName, userName)
	fs.SetupRepo(userName, userName, email)

	return templates.Page(
		"Home page",
		"", "",
		-1,
		nil,
		nil,
		nil,
	).Render(c.Request().Context(), c.Response())
}

func DemoPage(c echo.Context) error {
	return templates.DemoPage().Render(c.Request().Context(), c.Response())
}
