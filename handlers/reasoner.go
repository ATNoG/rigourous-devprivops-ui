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

func ReasonerMainPage(c echo.Context) error {
	reasonDir, err := fs.GetFile("reasoner")
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
		func() templ.Component { return templates.FileList("/reasoner/", "reasoner", ruleList) },
		nil,
		nil,
	).Render(c.Request().Context(), c.Response())
}

func ReasonerRuleEditor(c echo.Context) error {
	ruleName, err := url.QueryUnescape(c.Param("rule"))
	if err != nil {
		return err
	}

	ruleFileName := fmt.Sprintf("reasoner/%s", ruleName)
	ruleFile, err := fs.GetFile(ruleFileName)
	if err != nil {
		return err
	}

	ruleContent, err := os.ReadFile(ruleFile)
	if err != nil {
		return err
	}

	reasonDir, err := fs.GetFile("reasoner")
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
		func() templ.Component { return templates.FileList("/reasoner/", "reasoner/", ruleList) },
		func() templ.Component { return templates.EditorComponent("sparql", string(ruleContent), saveEndpoint) },
		nil,
	).Render(c.Request().Context(), c.Response())
}