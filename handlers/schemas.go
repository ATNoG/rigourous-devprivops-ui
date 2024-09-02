package handlers

import (
	"fmt"
	iofs "io/fs"
	"net/url"
	"os"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/Joao-Felisberto/devprivops-ui/util"
	"github.com/a-h/templ"
	"github.com/labstack/echo"
)

func SchemasMainPage(c echo.Context) error {
	userCookie, err := c.Cookie("username")
	if err != nil {
		return templates.Redirect("/").Render(c.Request().Context(), c.Response())
	}
	userName := userCookie.Value

	schemasDir, err := fs.GetFile("schemas", userName)
	if err != nil {
		return err
	}

	schemaFiles, err := os.ReadDir(schemasDir)
	if err != nil {
		return err
	}

	schemas := util.Map(schemaFiles, func(s iofs.DirEntry) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: s.Name(),
			Link: fmt.Sprintf("/schemas/%s", s.Name()),
		}
	})

	return templates.Page(
		"Schemas",
		"", "",
		templates.SCHEMAS,
		func() templ.Component { return templates.FileList("/schemas/", "schemas", schemas) },
		nil,
		nil,
	).Render(c.Request().Context(), c.Response())
}

func SchemaEditPage(c echo.Context) error {
	userCookie, err := c.Cookie("username")
	if err != nil {
		return templates.Redirect("/").Render(c.Request().Context(), c.Response())
	}
	userName := userCookie.Value

	schemaName, err := url.QueryUnescape(c.Param("schema"))
	if err != nil {
		return err
	}

	schemaFile, err := fs.GetFile(fmt.Sprintf("schemas/%s", schemaName), userName)
	if err != nil {
		return err
	}

	schemaContent, err := os.ReadFile(schemaFile)
	if err != nil {
		return err
	}

	schemasDir, err := fs.GetFile("schemas", userName)
	if err != nil {
		return err
	}

	schemaFiles, err := os.ReadDir(schemasDir)
	if err != nil {
		return err
	}

	schemas := util.Map(schemaFiles, func(s iofs.DirEntry) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: s.Name(),
			Link: fmt.Sprintf("/schemas/%s", s.Name()),
		}
	})

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape(schemaFile))
	return templates.Page(
		"Schemas",
		"schemaEditorContainer", "Visual",
		templates.SCHEMAS,
		func() templ.Component { return templates.FileList("/schemas/", "schemas", schemas) },
		func() templ.Component {
			return templates.SchemaEditor("yaml", string(schemaContent), saveEndpoint)
		},
		nil,
	).Render(c.Request().Context(), c.Response())
}
