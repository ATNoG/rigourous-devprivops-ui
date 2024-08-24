package handlers

import (
	"fmt"

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
			Link: "#",
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
				`config:
  - id: dpia:last update
<<<<<<< HEAD
    value: 02/03/2006
=======
    value: 02/03/2005
>>>>>>> 8837fc1f3f6cc392f7e690be6f4ef3cc44d8a6bf
  - id: dfd:as safeguards
    value:
      - encryption
  - id: dfd:as options
    value:
      - some option`,
				"#",
			)
		},
		nil,
	).Render(c.Request().Context(), c.Response())
}
