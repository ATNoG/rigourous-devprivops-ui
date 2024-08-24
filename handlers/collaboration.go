package handlers

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/Joao-Felisberto/devprivops-ui/util"
	"github.com/a-h/templ"
	"github.com/labstack/echo"
)

func Push(c echo.Context) error {
	cookie, err := c.Cookie("username")
	if err != nil {
		return err
	}
	userName := cookie.Value
	fmt.Printf("Username: %s\n", userName)

	repo := fmt.Sprintf("%s/%s", fs.LocalDir, userName)
	fmt.Printf("In repo '%s'\n", repo)
	conflictFiles, err := fs.GetConflicts(repo)
	if err != nil && err.Error() != "exit status 128" {
		fmt.Printf("Could not get conflicts: %s\n", err)
		return err
	}

	conflictList := util.Map(conflictFiles, func(file string) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: file,
			Link: fmt.Sprintf("/push/%s", url.QueryEscape(file)),
		}
	})
	// editor: https://github.com/microsoft/monaco-editor/issues/1529

	return templates.Page(
		"Merge",
		"", "",
		func() templ.Component { return templates.ConflictList(conflictList) },
		func() templ.Component {
			return templates.DiffEditor(
				"yaml",
				`aaa`,
				`bbb`,
				"#",
			)
		},
		nil,
	).Render(c.Request().Context(), c.Response())
}

func SolveMergeConflict(c echo.Context) error {
	cookie, err := c.Cookie("username")
	if err != nil {
		return err
	}
	userName := cookie.Value

	diffFile, err := url.QueryUnescape(c.Param("file"))
	if err != nil {
		return err
	}

	fmt.Printf("!!! %s: %s\n", diffFile, userName)
	modifiedFile, err := fs.GetFile(diffFile, userName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	originalFile, err := fs.GetFile(diffFile, "master")
	if err != nil {
		fmt.Println(err)
		return err
	}

	repo := fmt.Sprintf("%s/%s", fs.LocalDir, userName)
	fmt.Printf("In repo '%s'\n", repo)
	conflictFiles, err := fs.GetConflicts(repo)
	if err != nil && err.Error() != "exit status 128" {
		fmt.Printf("Could not get conflicts: %s\n", err)
		return err
	}

	conflictList := util.Map(conflictFiles, func(file string) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: file,
			Link: fmt.Sprintf("/push/%s", url.QueryEscape(file)),
		}
	})
	// editor: https://github.com/microsoft/monaco-editor/issues/1529

	originalContent, err := os.ReadFile(originalFile)
	if err != nil {
		return err
	}
	modifiedContent, err := os.ReadFile(modifiedFile)
	if err != nil {
		return err
	}

	fileKind := strings.Split(diffFile, "/")[0]

	var saveEndpoint string

	if strings.HasSuffix(fileKind, ".yml") || strings.HasSuffix(fileKind, ".yaml") {
		saveEndpoint = fmt.Sprintf("/save/%s", url.QueryEscape(diffFile))
	} else {
		switch fileKind {
		case "regulations":
			saveEndpoint = fmt.Sprintf("/save-regulation/%s", url.QueryEscape(diffFile))
		case "attack_trees":
			saveEndpoint = fmt.Sprintf("/save-tree/%s", url.QueryEscape(diffFile))
		case "report_data":
			saveEndpoint = "/save-report-data"
		case "requirements":
			saveEndpoint = "/save-requirements"
		default:
			saveEndpoint = fmt.Sprintf("/save/%s", url.QueryEscape(diffFile))
		}
	}

	return templates.Page(
		"Merge",
		"", "",
		func() templ.Component { return templates.ConflictList(conflictList) },
		func() templ.Component {
			return templates.DiffEditor(
				"yaml",
				string(originalContent),
				string(modifiedContent),
				saveEndpoint,
			)
		},
		nil,
	).Render(c.Request().Context(), c.Response())
}
