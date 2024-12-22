package handlers

import (
	"fmt"
	"net/http"

	sessionmanament "github.com/Joao-Felisberto/devprivops-ui/sessionManament"
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/Joao-Felisberto/devprivops-ui/tool"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/robert-nix/ansihtml"
)

// Endpoint to execute the analysis function of privguide on the user's repository
//
// `c`: the echo context
//
// returns: error if any internal function, like file reading, or template rendering fails.
func Analyse(c echo.Context) error {
	userName, err := sessionmanament.GetUsernameFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	fmt.Printf("User is '%s'\n", userName)
	res, err := tool.Analyse("", userName)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)

	htmlRes := ansihtml.ConvertToHTML([]byte(res))

	return templates.Page(
		"Analysis",
		"", "",
		templates.ANALYSE,
		nil,
		func() templ.Component { return templates.SimpleResult(string(htmlRes)) },
		nil,
	).Render(c.Request().Context(), c.Response())
}

// Endpoint to execute the test function of privguide on the user's repository
//
// `c`: the echo context
//
// returns: error if any internal function, like file reading, or template rendering fails.
func Test(c echo.Context) error {
	userName, err := sessionmanament.GetUsernameFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	res, err := tool.Test(userName)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)

	htmlRes := ansihtml.ConvertToHTML([]byte(res))

	return templates.Page(
		"Test",
		"", "",
		templates.TEST,
		nil,
		func() templ.Component { return templates.SimpleResult(string(htmlRes)) },
		nil,
	).Render(c.Request().Context(), c.Response())
}
