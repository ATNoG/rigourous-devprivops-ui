// This package agregates all handlers for the interface
package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
	sessionmanament "github.com/Joao-Felisberto/devprivops-ui/sessionManament"
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/Joao-Felisberto/devprivops-ui/util"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// Endpoint to manage conflicts visually.
//
// `c`: The echo context
//
// returns: error if any internal function, like file reading, or template rendering fails.
func MergeConflicts(c echo.Context) error {
	userName, err := sessionmanament.GetUsernameFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
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

	var rightBar func() templ.Component = nil
	if len(conflictFiles) == 0 {
		rightBar = func() templ.Component { return templates.ConflictSideBar() }
	}

	return templates.Page(
		"Merge",
		"", "",
		templates.CONFLICTS,
		func() templ.Component { return templates.ConflictList(conflictList) },
		nil,
		rightBar,
	).Render(c.Request().Context(), c.Response())
}

// Endpoint to manually solve conflicts on a file.
//
// `c`: The echo context
//
// # PATH PARAMETERS
//
// `file`: The name of the file with conflicts to solve
//
// returns: error if any internal function, like file reading, or template rendering fails.
func SolveMergeConflict(c echo.Context) error {
	userName, err := sessionmanament.GetUsernameFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	diffFile, err := url.QueryUnescape(c.Param("file"))
	if err != nil {
		return templates.Redirect("/").Render(c.Request().Context(), c.Response())
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

	if !util.Contains(conflictFiles, diffFile) {
		return templates.Redirect("/conflicts").Render(c.Request().Context(), c.Response())
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
		templates.CONFLICTS,
		func() templ.Component { return templates.ConflictList(conflictList) },
		func() templ.Component {
			return templates.DiffEditor(
				"yaml",
				string(originalContent),
				string(modifiedContent),
				saveEndpoint,
			)
		},
		// func() templ.Component { return templates.ConflictSideBar() },
		nil,
	).Render(c.Request().Context(), c.Response())
}

// Endpoint to push changes to the master repository, if there are no unsolved conflicts.
//
// `c`: The echo context
//
// returns: error if any internal function, like file reading, or template rendering fails.
func Push(c echo.Context) error {
	userName, err := sessionmanament.GetUsernameFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	// Do not push if there are still conflicts
	repo := fmt.Sprintf("%s/%s", fs.LocalDir, userName)
	conflictFiles, err := fs.GetConflicts(repo)
	if err != nil && err.Error() != "exit status 128" {
		fmt.Printf("Could not get conflicts: %s\n", err)
		return err
	}

	if len(conflictFiles) == 0 {

		fmt.Println("!!!!")
		out, err := fs.Push(userName)
		fmt.Println("!!!!")
		if err != nil && err.Error() != "exit status 128" {
			fmt.Println(out)
			fmt.Println(err)
			return err
		}
	} else {
		fmt.Println("There are still conflicts, please solve them first")
	}

	fmt.Println("Pushed to central repo")

	return nil
}
