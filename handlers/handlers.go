package handlers

import (
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/labstack/echo"
)

func Hello(c echo.Context) error {
	return templates.Hello("").Render(c.Request().Context(), c.Response())
}
