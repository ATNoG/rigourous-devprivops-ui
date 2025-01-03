package handlers

import (
	"fmt"
	iofs "io/fs"
	"net/http"
	"net/url"
	"os"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
	sessionmanament "github.com/Joao-Felisberto/devprivops-ui/sessionManament"
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/Joao-Felisberto/devprivops-ui/util"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// Endpoint to show all JSON schemas
//
// `c`: the echo context
//
// returns: error if any internal function, like file reading, or template rendering fails.
func SchemasMainPage(c echo.Context) error {
	userName, err := sessionmanament.GetUsernameFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
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

	return templates.Page(
		"Schemas",
		"", "",
		templates.SCHEMAS,
		func() templ.Component { return templates.FileList("/schemas/", "schemas", schemas) },
		nil,
		nil,
	).Render(c.Request().Context(), c.Response())
}

// Endpoint to edit a JSON schema
//
// `c`: the echo context
//
// # PATH PARAMETERS
//
// `schema`: The name of the schema file to edit
//
// returns: error if any internal function, like file reading, or template rendering fails.
func SchemaEditPage(c echo.Context) error {
	userName, err := sessionmanament.GetUsernameFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

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
