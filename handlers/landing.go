package handlers

import (
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// Endpoint to show all descriptions and the metadata file
//
// `c`: The echo context
//
// returns: error if any internal function, like file reading, or template rendering fails.
func LandingPage(c echo.Context) error {
	cookie, err := c.Cookie("username")
	if err != nil {
		return templates.Redirect("/").Render(c.Request().Context(), c.Response())
	}
	userName := cookie.Value

	// descs, err := fs.GetDescriptions("descriptions", userName)
	// if err != nil {
	// 	return err
	// }

	// descriptions := util.Map(descs, func(d string) templates.SideBarListElement {
	// 	return templates.SideBarListElement{
	// 		Text: d,
	// 		Link: fmt.Sprintf("/descriptions/%s", url.QueryEscape(d)),
	// 	}
	// })

	return templates.Page(
		"Home",
		"", "",
		templates.HOME,
		nil,
		func() templ.Component { return templates.HomeContent(userName) },
		nil,
	).Render(c.Request().Context(), c.Response())
}
