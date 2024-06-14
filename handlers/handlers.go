package handlers

import (
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/labstack/echo"
)

func DemoPage(c echo.Context) error {
	return templates.DemoPage().Render(c.Request().Context(), c.Response())
}
