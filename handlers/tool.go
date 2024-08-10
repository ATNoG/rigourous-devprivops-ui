package handlers

import (
	"fmt"

	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/Joao-Felisberto/devprivops-ui/tool"
	"github.com/a-h/templ"
	"github.com/labstack/echo"
	"github.com/robert-nix/ansihtml"
)

func Analyse(c echo.Context) error {
	res, err := tool.Analyse("")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)

	htmlRes := ansihtml.ConvertToHTML([]byte(res))

	return templates.Page(
		"Analysis",
		"", "",
		nil,
		func() templ.Component { return templates.SimpleResult(string(htmlRes)) },
		nil,
	).Render(c.Request().Context(), c.Response())
}

func Test(c echo.Context) error {
	res, err := tool.Test()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)

	htmlRes := ansihtml.ConvertToHTML([]byte(res))

	return templates.Page(
		"Test",
		"", "",
		nil,
		func() templ.Component { return templates.SimpleResult(string(htmlRes)) },
		nil,
	).Render(c.Request().Context(), c.Response())
}
