package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/Joao-Felisberto/devprivops-ui/util"
	"github.com/a-h/templ"
	"github.com/labstack/echo"
	"gopkg.in/yaml.v3"
)

func ExtraDataMainPage(c echo.Context) error {
	cookie, err := c.Cookie("username")
	if err != nil {
		return err
	}
	userName := cookie.Value

	extraDataFile, err := fs.GetFile("report_data/report_data.yml", userName)
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

	// saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape("report_data/report_data.yml"))
	saveEndpoint := "/save-report-data"
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
	cookie, err := c.Cookie("username")
	if err != nil {
		return err
	}
	userName := cookie.Value

	queryName, err := url.QueryUnescape(c.Param("query"))
	if err != nil {
		return err
	}

	queryFile, err := fs.GetFile(queryName, userName)
	if err != nil {
		return err
	}

	queryContent, err := os.ReadFile(queryFile)
	if err != nil {
		return err
	}

	extraDataFile, err := fs.GetFile("report_data/report_data.yml", userName)
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

	formLocation := c.FormValue("location")
	formHeading := c.FormValue("heading")
	formDescription := c.FormValue("description")
	formDataRowLine := c.FormValue("data row line")

	datum := util.First(contentList, func(d interface{}) bool {
		extraData := d.(map[string]interface{})
		file, err := fs.GetFile(extraData["query"].(string), userName)
		if err != nil {
			panic(err)
		}
		return file == queryFile
	})
	if datum == nil {
		fmt.Printf("Did not find corresponding extra data: '%s'\n", queryFile)
		return fmt.Errorf("did not find corresponding extra data: '%s'", queryFile)
	}
	extraDatum := (*datum).(map[string]interface{})

	if formLocation != "" || formHeading != "" || formDescription != "" || formDataRowLine != "" {
		if formLocation != "" {
			(*datum).(map[string]interface{})["location"] = formLocation
		}
		if formHeading != "" {
			(*datum).(map[string]interface{})["heading"] = formHeading
		}
		if formDescription != "" {
			(*datum).(map[string]interface{})["description"] = formDescription
		}
		if formDataRowLine != "" {
			(*datum).(map[string]interface{})["data row line"] = formDataRowLine
		}

		data, err := yaml.Marshal(datum)
		if err != nil {
			return err
		}

		extraDataFile, err := fs.GetFile("report_data/report_data.yml", userName)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("Writing to '%s'\n", extraDataFile)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = fs.WriteFileSync(extraDataFile, data, 0666)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	// saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape("report_data/report_data.yml"))
	saveEndpoint := "/save-report-data"
	return templates.Page(
		"Extra Data",
		"", "",
		func() templ.Component { return templates.FileList("/extra-data", "extra-data", extraDataList) },
		func() templ.Component {
			return templates.EditorComponent("sparql", string(queryContent), saveEndpoint)
		},
		func() templ.Component {
			return templates.SideBarForm(
				fmt.Sprintf("/extra-data/%s", c.Param("query")),
				templates.SideBarFormElement{
					Type:    templates.TEXT,
					Id:      "location",
					Label:   "Location",
					Default: extraDatum["location"].(string),
				},
				templates.SideBarFormElement{
					Type:    templates.TEXT,
					Id:      "heading",
					Label:   "Heading",
					Default: extraDatum["heading"].(string),
				},
				templates.SideBarFormElement{
					Type:    templates.TEXT,
					Id:      "description",
					Label:   "Description",
					Default: extraDatum["description"].(string),
				},
				templates.SideBarFormElement{
					Type:    templates.TEXT,
					Id:      "data row line",
					Label:   "Data row line",
					Default: extraDatum["data row line"].(string),
				},
			)
		},
	).Render(c.Request().Context(), c.Response())
}

func UpdateExtraData(c echo.Context) error {
	userCookie, err := c.Cookie("username")
	if err != nil {
		return err
	}
	userName := userCookie.Value

	emailCookie, err := c.Cookie("email")
	if err != nil {
		return err
	}
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
		fmt.Printf("failed to unmarshal: %s", err)
		return err
	}

	// sync files with the config
	dataDir, err := fs.GetFile("report_data/queries", userName)
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

	res, err := fs.AddAll(fs.LocalDir, userName)
	if err != nil {
		fmt.Println(res)
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	res, err = fs.Commit(fs.LocalDir, userName, strconv.Itoa(rand.Int()))
	if err != nil {
		fmt.Println(res)
		fmt.Println(err)
		return err
	}
	fmt.Println(res)

	return nil
}
