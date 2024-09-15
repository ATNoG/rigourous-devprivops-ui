package handlers

import (
	"fmt"
	"net/url"
	"os"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/Joao-Felisberto/devprivops-ui/util"

	"github.com/a-h/templ"
	"github.com/labstack/echo"
)

func TestOverview(c echo.Context) error {
	userCookie, err := c.Cookie("username")
	if err != nil {
		return templates.Redirect("/").Render(c.Request().Context(), c.Response())
	}
	userName := userCookie.Value

	testDirs, err := fs.GetTests(userName)
	if err != nil {
		return err
	}

	testScenarios := util.Map(testDirs, func(dir string) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: dir,
			Link: fmt.Sprintf("/tests/%s", dir),
		}
	})

	return templates.Page(
		"Regulations",
		"regulation-editor", "Visual",
		templates.REGULATIONS,
		func() templ.Component { return templates.RegulationList("tests", testScenarios) },
		nil,
		nil,
	).Render(c.Request().Context(), c.Response())
}

func TestScenarioSelect(c echo.Context) error {
	userCookie, err := c.Cookie("username")
	if err != nil {
		return templates.Redirect("/").Render(c.Request().Context(), c.Response())
	}
	userName := userCookie.Value

	scenarioName := c.Param("scenario")
	fmt.Printf("Scenario '%s'\n", scenarioName)

	testDirs, err := fs.GetTests(userName)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return err
	}

	testScenarios := util.Map(testDirs, func(dir string) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: dir,
			Link: fmt.Sprintf("/tests/%s", dir),
		}
	})

	scenarioDescFiles, err := fs.GetTestScenarios(scenarioName, userName)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return err
	}
	fmt.Printf("Found %d scenarios", len(scenarioDescFiles))

	scenarioDescs := util.Map(scenarioDescFiles, func(dir string) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: dir,
			Link: fmt.Sprintf("/tests/%s/%s", scenarioName, url.QueryEscape(dir)),
		}
	})

	return templates.Page(
		"Regulations",
		"regulation-editor", "Visual",
		templates.REGULATIONS,
		func() templ.Component {
			return templates.VerticalList(
				func() templ.Component { return templates.RegulationList("tests", testScenarios) },
				func() templ.Component { return templates.SideBarList(scenarioDescs) },
			)
		},
		/*
			func() templ.Component {
				return templates.RegulationEditor("yaml", string(cfgContent), saveEndpoint, &jsonString)
			},
		*/
		nil,
		nil,
	).Render(c.Request().Context(), c.Response())
}

func TestScenarioEdit(c echo.Context) error {
	userCookie, err := c.Cookie("username")
	if err != nil {
		fmt.Println("ERROR: ", err)
		return templates.Redirect("/").Render(c.Request().Context(), c.Response())
	}
	userName := userCookie.Value

	scenarioName := c.Param("scenario")

	descFile, err := url.QueryUnescape(c.Param("desc"))
	if err != nil {
		fmt.Println("ERROR: ", err)
		return err
	}

	testDirs, err := fs.GetTests(userName)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return err
	}

	testScenarios := util.Map(testDirs, func(dir string) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: dir,
			Link: fmt.Sprintf("/tests/%s", dir),
		}
	})

	scenarioDescFiles, err := fs.GetTestScenarios(scenarioName, userName)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return err
	}

	scenarioDescs := util.Map(scenarioDescFiles, func(dir string) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: dir,
			Link: fmt.Sprintf("/tests/%s/%s", scenarioName, url.QueryEscape(dir)),
		}
	})

	fmt.Printf(">>> ! %s", descFile)
	descPath, err := fs.GetFile(descFile, userName)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return err
	}
	cfgContent, err := os.ReadFile(descPath)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return err
	}

	return templates.Page(
		"Regulations",
		"regulation-editor", "Visual",
		templates.REGULATIONS,
		func() templ.Component {
			return templates.VerticalList(
				func() templ.Component { return templates.RegulationList("tests", testScenarios) },
				func() templ.Component { return templates.SideBarList(scenarioDescs) },
			)
		},
		func() templ.Component {
			return templates.EditorComponent("yaml", string(cfgContent), "/")
		},
		nil,
	).Render(c.Request().Context(), c.Response())
}
