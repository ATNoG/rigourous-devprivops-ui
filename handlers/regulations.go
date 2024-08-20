package handlers

import (
	"encoding/json"
	"fmt"
	"io"
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

func RegulationsMainPage(c echo.Context) error {
	userCookie, err := c.Cookie("username")
	if err != nil {
		return err
	}
	userName := userCookie.Value

	regulationDirs, err := fs.GetRegulations(userName)
	if err != nil {
		return err
	}

	regulations := util.Map(regulationDirs, func(r string) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: r,
			Link: fmt.Sprintf("/regulations/%s", r),
		}
	})

	return templates.Page(
		"Regulations",
		"", "",
		func() templ.Component {
			return templates.RegulationList("regulations", regulations)
		},
		nil,
		nil,
	).Render(c.Request().Context(), c.Response())
}

func RegulationView(c echo.Context) error {
	userCookie, err := c.Cookie("username")
	if err != nil {
		return err
	}
	userName := userCookie.Value

	regName := c.Param("reg")

	regCfgFileName := fmt.Sprintf("regulations/%s/policies.yml", regName)
	regulationFile, err := fs.GetFile(regCfgFileName, userName)
	if err != nil {
		return err
	}

	cfgContent, err := os.ReadFile(regulationFile)
	if err != nil {
		return err
	}

	regulationDirs, err := fs.GetRegulations(userName)
	if err != nil {
		return err
	}

	regulations := util.Map(regulationDirs, func(r string) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: r,
			Link: fmt.Sprintf("/regulations/%s", r),
		}
	})

	var policies []interface{} = []interface{}{}

	err = yaml.Unmarshal(cfgContent, &policies)
	if err != nil {
		return err
	}

	jsonContent, err := json.Marshal(&policies)
	if err != nil {
		return err
	}
	jsonString := string(jsonContent)

	policyFiles := util.Map(policies, func(pol interface{}) templates.SideBarListElement {
		p := pol.(map[string]interface{})
		fmt.Printf("!! /policy/%s\n", url.QueryEscape(p["file"].(string)))
		return templates.SideBarListElement{
			Text: p["file"].(string),
			Link: fmt.Sprintf("/policy/%s", url.QueryEscape(p["file"].(string))),
		}
	})

	// regulations = append(regulations, policyFiles...)

	// saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape(regCfgFileName))
	saveEndpoint := fmt.Sprintf("/save-regulation/%s", url.QueryEscape(regCfgFileName))

	// TODO: tests
	return templates.Page(
		"Regulations",
		"regulation-editor", "Visual",
		func() templ.Component {
			return templates.VerticalList(
				func() templ.Component { return templates.RegulationList("regulations", regulations) },
				func() templ.Component { return templates.SideBarList(policyFiles) },
			)
			// return templates.SideBarList(regulations)
		},
		// func() templ.Component { return templates.EditorComponent("yaml", string(cfgContent), saveEndpoint) },
		func() templ.Component {
			return templates.RegulationEditor("yaml", string(cfgContent), saveEndpoint, &jsonString)
		},
		nil,
	).Render(c.Request().Context(), c.Response())
}

func PolicyEdit(c echo.Context) error {
	userCookie, err := c.Cookie("username")
	if err != nil {
		return err
	}
	userName := userCookie.Value

	polName, err := url.QueryUnescape(c.Param("pol"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	regulationName := strings.Split(polName, "/")[1]
	regCfgFileName := fmt.Sprintf("regulations/%s/policies.yml", regulationName)
	regCfgFile, err := fs.GetFile(regCfgFileName, userName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	regCfgContent, err := os.ReadFile(regCfgFile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	polFile, err := fs.GetFile(polName, userName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	polContent, err := os.ReadFile(polFile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	regulationDirs, err := fs.GetRegulations(userName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	regulations := util.Map(regulationDirs, func(r string) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: r,
			Link: fmt.Sprintf("/regulations/%s", r),
		}
	})

	var policies []interface{} = []interface{}{}

	err = yaml.Unmarshal(regCfgContent, &policies)
	if err != nil {
		fmt.Println(string(regCfgContent))
		fmt.Println(err)
		return err
	}

	policyFiles := util.Map(policies, func(pol interface{}) templates.SideBarListElement {
		p := pol.(map[string]interface{})
		return templates.SideBarListElement{
			Text: p["file"].(string),
			Link: p["file"].(string),
		}
	})

	regulations = append(regulations, policyFiles...)

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape(polFile))

	formTitle := c.FormValue("title")
	formDescription := c.FormValue("description")
	formConsistency := c.FormValue("consistency") // TODO: can't turn off consistency once turned on
	formViolations := c.FormValue("violations")
	formMapping := c.FormValue("mapping")

	policyRaw := util.First(policies, func(pol interface{}) bool {
		p := pol.(map[string]interface{})
		return p["file"] == polName
	})

	if formTitle != "" || formDescription != "" || formConsistency != "" || formViolations != "" || formMapping != "" {
		fmt.Printf("-> '%s' '%s' '%s' '%s' '%s'\n",
			formTitle,
			formDescription,
			formConsistency, // TODO: can't turn off consistency once turned on
			formViolations,
			formMapping,
		)

		if policyRaw == nil {
			fmt.Printf("ERROR No policy for '%s' found\n", polName)
			return nil
		}

		if formTitle != "" {
			((*policyRaw).(map[string]interface{}))["title"] = formTitle
		}
		if formDescription != "" {
			((*policyRaw).(map[string]interface{}))["description"] = formDescription
		}
		if formConsistency != "" {
			((*policyRaw).(map[string]interface{}))["is consistency"] = formConsistency
		}
		if formViolations != "" {
			((*policyRaw).(map[string]interface{}))["maximum violations"] = formViolations
		}
		if formMapping != "" {
			((*policyRaw).(map[string]interface{}))["mapping message"] = formMapping
		}

		data, err := yaml.Marshal(policies)
		if err != nil {
			return err
		}

		regulationFile, err := fs.GetFile(fmt.Sprintf("regulations/%s/policies.yml", regulationName), userName)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("Writing to '%s'\n", regulationFile)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = fs.WriteFileSync(regulationFile, data, 0666)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	policy := (*policyRaw).(map[string]interface{})

	return templates.Page(
		"Regulations",
		"", "",
		func() templ.Component {
			return templates.SideBarList(regulations)
		},
		func() templ.Component { return templates.EditorComponent("sparql", string(polContent), saveEndpoint) },
		func() templ.Component {
			return templates.SideBarForm(fmt.Sprintf("/policy/%s", c.Param("pol")),
				templates.SideBarFormElement{
					Type:    templates.TEXT,
					Id:      "title",
					Label:   "Title",
					Default: policy["title"].(string),
				},
				templates.SideBarFormElement{
					Type:    templates.TEXT,
					Id:      "description",
					Label:   "Description",
					Default: policy["description"].(string),
				},
				templates.SideBarFormElement{
					Type:    templates.CHECKBOX,
					Id:      "consistency",
					Label:   "Is consistency",
					Default: strconv.FormatBool(policy["is consistency"].(bool)),
				},
				templates.SideBarFormElement{
					Type:    templates.TEXT,
					Id:      "violations",
					Label:   "Maximum violations",
					Default: fmt.Sprintf("%d", policy["maximum violations"].(int)),
				},
				templates.SideBarFormElement{
					Type:    templates.TEXT,
					Id:      "mapping",
					Label:   "Mapping message",
					Default: policy["mapping message"].(string),
				},
			)
		},
	).Render(c.Request().Context(), c.Response())
}

func CreateRegulation(c echo.Context) error {
	file := c.QueryParam("path")
	path := fmt.Sprintf("%s/%s", fs.LocalDir, file)

	fmt.Printf("Create REGULATION '%s'\n", path)

	if err := os.Mkdir(path, 0777); err != nil {
		fmt.Println(err)
		return err
	}
	if err := os.Mkdir(fmt.Sprintf("%s/consistency", path), 0755); err != nil {
		fmt.Println(err)
		return err
	}
	if err := os.Mkdir(fmt.Sprintf("%s/policies", path), 0755); err != nil {
		fmt.Println(err)
		return err
	}

	if err := fs.WriteFileSync(fmt.Sprintf("%s/policies.yml", path), []byte("[]"), 0666); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func DeleteRegulation(c echo.Context) error {
	file := c.QueryParam("path")
	path := fmt.Sprintf("%s/regulations/%s", fs.LocalDir, file)

	fmt.Printf("Delete REGULATION '%s'\n", path)

	if err := os.RemoveAll(path); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func UpdateRegulation(c echo.Context) error {
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

	fName, err := url.QueryUnescape(c.Param("reg"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	contents := []interface{}{}
	json.Unmarshal(body, &contents)

	// sync files with the config
	pathComponents := strings.Split(fName, "/")
	regPath := strings.Join(pathComponents[:len(pathComponents)-1], "/") + "/consistency"
	realRegPath, err := fs.GetFile(regPath, userName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Syncing files at %s\n", realRegPath)
	err = util.SyncFileList(realRegPath, util.Map(contents, func(e interface{}) string {
		m := e.(map[string]interface{})
		components := strings.Split(m["file"].(string), "/")
		return components[len(components)-1]
	}))
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

	fmt.Println(string(data))

	file, err := fs.GetFile(fName, userName)
	if err != nil {
		return err
	}

	fmt.Printf("Writing to %s: %s \n", file, string(data))

	if err := fs.WriteFileSync(file, data, 0666); err != nil {
		return err
	}

	fmt.Printf("In %s: %s %s\n", fs.LocalDir, userName, email)

	return nil
}
