package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/Joao-Felisberto/devprivops-ui/util"
	"github.com/a-h/templ"
	"github.com/labstack/echo"
	"gopkg.in/yaml.v3"
)

func ExtraDataMainPage(c echo.Context) error {
	extraDataFile, err := fs.GetFile("report_data/report_data.yml")
	if err != nil {
		return err
	}

	extraDataContent, err := os.ReadFile(extraDataFile)
	if err != nil {
		return err
	}

	contentList := []interface{}{}
	err = yaml.Unmarshal(extraDataContent, &contentList)
	if err != nil {
		return err
	}

	rawJsonData, err := json.Marshal(contentList)
	if err != nil {
		return err
	}
	jsonData := string(rawJsonData)

	extraDataList := util.Map(contentList, func(ed interface{}) templates.SideBarListElement {
		content := ed.(map[string]interface{})

		return templates.SideBarListElement{
			Text: content["query"].(string),
			Link: fmt.Sprintf("/extra-data/%s", url.QueryEscape(content["query"].(string))),
		}
	})

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape("report_data/report_data.yml"))
	return templates.Page(
		"Extra Data",
		"extra-data-editor", "Visual",
		func() templ.Component { return templates.FileList("/extra-data", "extra-data", extraDataList) },
		func() templ.Component {
			// return templates.EditorComponent("yaml", string(extraDataContent), saveEndpoint)
			return templates.ExtraDataEditor("yaml", string(extraDataContent), saveEndpoint, &jsonData)
		},
		nil,
	).Render(c.Request().Context(), c.Response())
}

func ExtraDataQuery(c echo.Context) error {
	queryName, err := url.QueryUnescape(c.Param("query"))
	if err != nil {
		return err
	}

	queryFile, err := fs.GetFile(queryName)
	if err != nil {
		return err
	}

	queryContent, err := os.ReadFile(queryFile)
	if err != nil {
		return err
	}

	extraDataFile, err := fs.GetFile("report_data/report_data.yml")
	if err != nil {
		return err
	}

	extraDataContent, err := os.ReadFile(extraDataFile)
	if err != nil {
		return err
	}

	contentList := []interface{}{}
	err = yaml.Unmarshal(extraDataContent, &contentList)
	if err != nil {
		return err
	}

	extraDataList := util.Map(contentList, func(ed interface{}) templates.SideBarListElement {
		content := ed.(map[string]interface{})

		return templates.SideBarListElement{
			Text: content["query"].(string),
			Link: fmt.Sprintf("/extra-data/%s", url.QueryEscape(content["query"].(string))),
		}
	})

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape("report_data/report_data.yml"))
	return templates.Page(
		"Extra Data",
		"", "",
		func() templ.Component { return templates.FileList("/extra-data", "extra-data", extraDataList) },
		func() templ.Component {
			return templates.EditorComponent("sparql", string(queryContent), saveEndpoint)
		},
		nil,
	).Render(c.Request().Context(), c.Response())
}

func UpdateExtraData(c echo.Context) error {
	userCookie := util.Filter(c.Request().Cookies(), func(cookie *http.Cookie) bool {
		return cookie.Name == "username"
	})[0]
	emailCookie := util.Filter(c.Request().Cookies(), func(cookie *http.Cookie) bool {
		return cookie.Name == "email"
	})[0]

	userName := userCookie.Value
	email := emailCookie.Value

	/*
		fName, err := url.QueryUnescape(c.Param("tree"))
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("Save to: %s\n", fName)
	*/

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
	dataDir, err := fs.GetFile("report_data/queries")
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Syncing files at %s\n", dataDir)
	err = util.SyncFileList(dataDir, util.Map(contents, func(e interface{}) string {
		m := e.(map[string]interface{})
		components := strings.Split(m["query"].(string), "/")
		return components[len(components)-1]
	}))
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Files synced")

	data, err := yaml.Marshal(contents)
	// _, err = yaml.Marshal(contents)
	// fmt.Printf("%+v\n", contents)

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Data to write \\/")
	fmt.Println(string(data))

	/*
		file, err := fs.GetFile(fName)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Printf("Writing to %s: %s \n", file, string(data))

		if err := os.WriteFile(file, data, 0666); err != nil {
			return err
		}
	*/

	fmt.Printf("In %s: %s %s\n", fs.LocalDir, userName, email)

	return nil
}
