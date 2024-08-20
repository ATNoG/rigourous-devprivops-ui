package handlers

import (
	"net/http"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/a-h/templ"
	"github.com/labstack/echo"
)

func HomePage(c echo.Context) error {
	return templates.Page(
		"Home page",
		"", "",
		nil,
		func() templ.Component { return templates.LoginForm() },
		nil,
	).Render(c.Request().Context(), c.Response())
}

func LogIn(c echo.Context) error {
	userNameCookie := new(http.Cookie)
	userNameCookie.Name = "username"
	userNameCookie.Value = c.FormValue("username")
	userNameCookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(userNameCookie)

	mailCookie := new(http.Cookie)
	mailCookie.Name = "email"
	mailCookie.Value = c.FormValue("email")
	mailCookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(mailCookie)

	fs.SessionManager.AddSession(c.FormValue("username"), "master")

	return templates.Page(
		"Home page",
		"", "",
		nil,
		func() templ.Component { return templates.LoginForm() },
		nil,
	).Render(c.Request().Context(), c.Response())
}

func DemoPage(c echo.Context) error {
	return templates.DemoPage().Render(c.Request().Context(), c.Response())
}
