package handlers

import (
	"fmt"
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

// Endpoint to view all tests and test specifications
//
// `c`: the echo context
//
// returns: error if any internal function, like file reading, or template rendering fails.
func TestOverview(c echo.Context) error {
	userName, err := sessionmanament.GetUsernameFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

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

	testSpecsFile, err := fs.GetFile("tests/spec.json", userName)
	if err != nil {
		return err
	}

	testSpecs, err := os.ReadFile(testSpecsFile)
	if err != nil {
		return err
	}

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape("tests/spec.json"))
	return templates.Page(
		"Tests",
		"test-editor", "Visual",
		templates.TEST,
		func() templ.Component { return templates.RegulationList("tests", testScenarios) },
		func() templ.Component {
			// return templates.EditorComponent("json", string(testSpecs), saveEndpoint)
			return templates.TestEditor("json", string(testSpecs), saveEndpoint)
		},
		nil,
	).Render(c.Request().Context(), c.Response())
}

// Endpoint to view all test descriptions for a regulation
//
// `c`: the echo context
//
// # PATH PARAMETERS
//
// `scenario`: The name of the scenario to edit
//
// returns: error if any internal function, like file reading, or template rendering fails.
func TestScenarioSelect(c echo.Context) error {
	userName, err := sessionmanament.GetUsernameFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

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
		"Tests",
		"", "",
		templates.TEST,
		func() templ.Component {
			return templates.VerticalList(
				func() templ.Component { return templates.RegulationList("tests", testScenarios) },
				func() templ.Component { return templates.SideBarList(scenarioDescs) },
			)
		},
		nil,
		nil,
	).Render(c.Request().Context(), c.Response())
}

// Endpoint to edit a test scenario
//
// `c`: the echo context
//
// # PATH PARAMETERS
//
// `scenario`: The name of the scenario to edit
//
// `desc`: The path to the description file
//
// returns: error if any internal function, like file reading, or template rendering fails.
func TestScenarioEdit(c echo.Context) error {
	userName, err := sessionmanament.GetUsernameFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

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

	fmt.Printf(">>> ! %s\n", descFile)
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

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape(descPath))
	return templates.Page(
		"Tests",
		"", "",
		templates.TEST,
		func() templ.Component {
			return templates.VerticalList(
				func() templ.Component { return templates.RegulationList("tests", testScenarios) },
				func() templ.Component { return templates.SideBarList(scenarioDescs) },
			)
		},
		func() templ.Component {
			return templates.EditorWithVisualizer("yaml", string(cfgContent), saveEndpoint)
		},
		nil,
	).Render(c.Request().Context(), c.Response())
}

/*
func SaveTestSpecs(c echo.Context) error {
	userCookie, err := c.Cookie("username")
	if err != nil {
		return templates.Redirect("/").Render(c.Request().Context(), c.Response())
	}
	userName := userCookie.Value

	emailCookie, err := c.Cookie("email")
	if err != nil {
		return templates.Redirect("/").Render(c.Request().Context(), c.Response())
	}
	email := emailCookie.Value

	body, err := io.ReadAll(c.Request().Body)

	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))

	var contents []interface{}
	err = json.Unmarshal(body, &contents)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// sync files with the config
	reqPath, err := fs.GetFile("", userName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Syncing files at %s\n", reqPath)
	fileList := util.Flatten(util.Map(contents, func(e interface{}) []string {
		useCase := e.(map[string]interface{})
		requirements := useCase["requirements"].([]interface{})
		queryFiles := util.Map(requirements, func(r interface{}) string {
			req := r.(map[string]interface{})
			query := req["query"].(string)
			components := strings.Split(query, "/")
			return components[len(components)-1]
		})

		return queryFiles
	}))
	fileList = append(fileList, "requirements.yml")
	err = util.SyncFileList(reqPath, fileList)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Files synced")

	data, err := yaml.Marshal(contents)

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Data to write \\/")
	fmt.Println(string(data))

	file, err := fs.GetFile("requirements/requirements.yml", userName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Writing to %s: %s \n", file, string(data))

	if err := fs.WriteFileSync(file, data, 0666); err != nil {
		return err
	}

	fmt.Printf("In %s: %s %s\n", fs.LocalDir, userName, email)

	res, err := fs.AddAll(userName)
	if err != nil {
		fmt.Println(res)
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	res, err = fs.Commit(userName, strconv.Itoa(rand.Int()))
	if err != nil {
		fmt.Println(res)
		fmt.Println(err)
		return err
	}
	fmt.Println(res)

	return nil
}
*/
