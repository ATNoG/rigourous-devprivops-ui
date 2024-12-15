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
	"github.com/labstack/echo/v4"
)

// Endpoint to show all reasoner rules
//
// `c`: The echo context
//
// returns: error if any internal function, like file reading, or template rendering fails.
func ReasonerMainPage(c echo.Context) error {
	userCookie, err := c.Cookie("username")
	if err != nil {
		return templates.Redirect("/").Render(c.Request().Context(), c.Response())
	}
	userName := userCookie.Value

	reasonDir, err := fs.GetFile("reasoner", userName)
	if err != nil {
		return err
	}
	ruleFiles, err := os.ReadDir(reasonDir)
	if err != nil {
		return err
	}

	ruleList := util.Map(ruleFiles, func(t iofs.DirEntry) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: t.Name(),
			Link: fmt.Sprintf("/reasoner/%s", url.QueryEscape(t.Name())),
		}
	})

	return templates.Page(
		"Reasoner",
		"", "",
		templates.REASONER,
		func() templ.Component { return templates.FileList("/reasoner/", "reasoner", ruleList) },
		nil,
		nil,
	).Render(c.Request().Context(), c.Response())
}

// Endpoint to edit a single reasoner rule
//
// `c`: The echo context
//
// # PATH PARAMETERS
//
// `rule`: The reasoner rule file to edit
//
// returns: error if any internal function, like file reading, or template rendering fails.
func ReasonerRuleEditor(c echo.Context) error {
	userCookie, err := c.Cookie("username")
	if err != nil {
		return templates.Redirect("/").Render(c.Request().Context(), c.Response())
	}
	userName := userCookie.Value

	ruleName, err := url.QueryUnescape(c.Param("rule"))
	if err != nil {
		return err
	}

	ruleFileName := fmt.Sprintf("reasoner/%s", ruleName)
	ruleFile, err := fs.GetFile(ruleFileName, userName)
	if err != nil {
		return err
	}

	ruleContent, err := os.ReadFile(ruleFile)
	if err != nil {
		return err
	}

	reasonDir, err := fs.GetFile("reasoner", userName)
	if err != nil {
		return err
	}
	ruleFiles, err := os.ReadDir(reasonDir)
	if err != nil {
		return err
	}

	ruleList := util.Map(ruleFiles, func(t iofs.DirEntry) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: t.Name(),
			Link: fmt.Sprintf("/reasoner/%s", url.QueryEscape(t.Name())),
		}
	})

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape(ruleFileName))

	return templates.Page(
		"Reasoner",
		"", "",
		templates.REASONER,
		func() templ.Component { return templates.FileList("/reasoner/", "reasoner/", ruleList) },
		func() templ.Component { return templates.EditorComponent("sparql", string(ruleContent), saveEndpoint) },
		nil,
	).Render(c.Request().Context(), c.Response())
}
