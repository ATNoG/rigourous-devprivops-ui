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
	"gopkg.in/yaml.v3"
)

func DescriptionsMainPage(c echo.Context) error {
	descs, err := fs.GetDescriptions("descriptions")
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
		func() templ.Component { return templates.FileList("descriptions", "descriptions", descriptions) },
		nil,
		nil,
	).Render(c.Request().Context(), c.Response())
}

func DescriptionEdit(c echo.Context) error {
	cookie, err := c.Cookie("username")
	if err != nil {
		return err
	}

	fmt.Printf("User: %s\n", cookie.Value)

	desc, err := url.QueryUnescape(c.Param("desc"))
	if err != nil {
		return err
	}

	descFile, err := fs.GetFile(desc)
	if err != nil {
		return err
	}

	descContent, err := os.ReadFile(descFile)
	if err != nil {
		return err
	}

	descs, err := fs.GetDescriptions("descriptions")
	if err != nil {
		return err
	}

	urisFile, err := fs.GetFile("uris.yml")
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

	return templates.Page(
		"My page",
		"graphContainer", "Visual",
		func() templ.Component { return templates.FileList("descriptions", "descriptions", descriptions) },
		//		func() templ.Component { return templates.EditorComponent("yaml", string(descContent), saveEndpoint) },
		func() templ.Component {
			return templates.EditorWithVisualizer("yaml", string(descContent), saveEndpoint)
		},
		func() templ.Component {
			return templates.DescriptionMetadata(
				fmt.Sprintf("/descriptions/%s", c.Param("desc")),
				uri,
				uriList,
			)
		},
		// nil,
	).Render(c.Request().Context(), c.Response())
}
