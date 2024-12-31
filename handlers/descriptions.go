package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
	sessionmanament "github.com/Joao-Felisberto/devprivops-ui/sessionManament"
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/Joao-Felisberto/devprivops-ui/util"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"
)

var STATIC_DIR = ""

// Endpoint to show all descriptions and the metadata file
//
// `c`: The echo context
//
// returns: error if any internal function, like file reading, or template rendering fails.
func DescriptionsMainPage(c echo.Context) error {
	userName, err := sessionmanament.GetUsernameFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	descs, err := fs.GetDescriptions("descriptions", userName)
	if err != nil {
		return err
	}

	descriptions := util.Map(descs, func(d string) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: d,
			Link: fmt.Sprintf("/descriptions/%s", url.QueryEscape(d)),
		}
	})

	return templates.Page(
		"My page",
		"", "",
		templates.DESCRIPTIONS,
		func() templ.Component { return templates.FileList("descriptions", "descriptions", descriptions) },
		nil,
		nil,
	).Render(c.Request().Context(), c.Response())
}

// Endpoint to edit a single description file
//
// `c`: The echo context
//
// # PATH PARAMETERS
//
// `desc`: The description file to edit
//
// returns: error if any internal function, like file reading, or template rendering fails.
func DescriptionEdit(c echo.Context) error {
	userName, err := sessionmanament.GetUsernameFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	fmt.Printf("User: %s\n", userName)

	desc, err := url.QueryUnescape(c.Param("desc"))
	if err != nil {
		return err
	}

	descFile, err := fs.GetFile(desc, userName)
	if err != nil {
		return err
	}

	descContent, err := os.ReadFile(descFile)
	if err != nil {
		return err
	}

	descs, err := fs.GetDescriptions("descriptions", userName)
	if err != nil {
		return err
	}

	urisFile, err := fs.GetFile("uris.yml", userName)
	if err != nil {
		return err
	}
	urisContent, err := os.ReadFile(urisFile)
	if err != nil {
		return err
	}

	uris := []map[string]interface{}{}
	err = yaml.Unmarshal(urisContent, &uris)
	if err != nil {
		return err
	}

	uriList := util.Map(uris, func(uri map[string]interface{}) string {
		return uri["abreviation"].(string)
	})

	uri := c.FormValue("uri")
	// fmt.Printf("Old %s\n", uri)
	if uri != "" {
		for _, e := range uris {
			if e["abreviation"].(string) != uri {
				continue
			}

			files := e["files"].([]interface{})
			if util.First(files, func(f interface{}) bool {
				return f.(string) == desc
			}) == nil {
				e["files"] = append(e["files"].([]interface{}), desc)
			}
		}

		newContent, err := yaml.Marshal(uris)
		if err != nil {
			return err
		}
		err = fs.WriteFileSync(urisFile, newContent, 0666)
		if err != nil {
			return err
		}
	}

	newURIAbrev := c.FormValue("abreviation")
	newURI := c.FormValue("new-uri")
	if newURIAbrev != "" && newURI != "" {
		fmt.Printf("%s[%s]\n", newURIAbrev, newURI)

		for _, e := range uris {
			if e["abreviation"].(string) != uri {
				e["files"] = util.Filter(e["files"].([]interface{}), func(f interface{}) bool {
					return f.(string) != desc
				})
			}
		}

		uris = append(uris, map[string]interface{}{
			"abreviation": newURIAbrev,
			"uri":         newURI,
			"files":       []string{desc},
		})

		newContent, err := yaml.Marshal(uris)
		if err != nil {
			return err
		}
		err = fs.WriteFileSync(urisFile, newContent, 0666)
		if err != nil {
			return err
		}
	}

	descriptions := util.Map(descs, func(d string) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: d,
			Link: fmt.Sprintf("/descriptions/%s", url.QueryEscape(d)),
		}
	})

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape(desc))

	actualURI := util.First(uris, func(metadata map[string]interface{}) bool {
		files := util.Map(metadata["files"].([]interface{}), func(e interface{}) string {
			fmt.Println(e.(string))
			return e.(string)
		})
		return slices.ContainsFunc(files, func(f string) bool {
			matched, err := regexp.Match(f, []byte(desc))
			if err != nil {
				panic(err)
			}

			return matched
		})
	})

	tmp := strings.Split(descFile, ".")
	pluginName := tmp[len(tmp)-2]
	pluginPath := fmt.Sprintf("%s/%s/%s.js", STATIC_DIR, pluginName, pluginName)

	if _, err := os.Stat(pluginPath); errors.Is(err, os.ErrNotExist) {
		pluginPath = ""
	} else if !strings.HasPrefix(pluginPath, "/") {
		pluginPath = fmt.Sprintf("/%s", pluginPath)
	}

	return templates.Page(
		"My page",
		"graphContainer", "Visual",
		templates.DESCRIPTIONS,
		func() templ.Component { return templates.FileList("descriptions", "descriptions", descriptions) },
		//		func() templ.Component { return templates.EditorComponent("yaml", string(descContent), saveEndpoint) },
		func() templ.Component {
			return templates.EditorWithVisualizer("yaml", string(descContent), saveEndpoint, pluginPath)
		},
		func() templ.Component {
			return templates.DescriptionMetadata(
				fmt.Sprintf("/descriptions/%s", c.Param("desc")),
				(*actualURI)["abreviation"].(string),
				(*actualURI)["uri"].(string),
				uriList,
			)
		},
		// nil,
	).Render(c.Request().Context(), c.Response())
}
