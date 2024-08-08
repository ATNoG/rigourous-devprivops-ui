package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	iofs "io/fs"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/Joao-Felisberto/devprivops-ui/tool"
	"github.com/Joao-Felisberto/devprivops-ui/util"
	"github.com/a-h/templ"
	"github.com/labstack/echo"
	"gopkg.in/yaml.v3"
)

func HomePage(c echo.Context) error {
	return templates.Page(
		"Home page",
		"", "",
		nil,
		func() templ.Component { return templates.LoginForm() },
		nil,
	).Render(c.Request().Context(), c.Response())
}

func LogIn(c echo.Context) error {
	userNameCookie := new(http.Cookie)
	userNameCookie.Name = "username"
	userNameCookie.Value = c.FormValue("username")
	userNameCookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(userNameCookie)

	mailCookie := new(http.Cookie)
	mailCookie.Name = "email"
	mailCookie.Value = c.FormValue("email")
	mailCookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(mailCookie)

	return templates.Page(
		"Home page",
		"", "",
		nil,
		func() templ.Component { return templates.LoginForm() },
		nil,
	).Render(c.Request().Context(), c.Response())
}

func DemoPage(c echo.Context) error {
	return templates.DemoPage().Render(c.Request().Context(), c.Response())
}

func TreesMainPage(c echo.Context) error {
	atkDir, err := fs.GetFile("attack_trees/descriptions/")
	if err != nil {
		fmt.Println(err)
		return err
	}
	treeFiles, err := os.ReadDir(atkDir)
	if err != nil {
		fmt.Println(err)
		return err
	}

	treeList := util.Map(treeFiles, func(t iofs.DirEntry) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: t.Name(),
			Link: fmt.Sprintf("/trees/%s", url.QueryEscape(t.Name())),
		}
	})

	return templates.Page(
		"Trees",
		"", "",
		func() templ.Component {
			return templates.FileList("/trees/", "attack_trees/descriptions", treeList)
		},
		nil,
		nil,
	).Render(c.Request().Context(), c.Response())
}

func TreeView(c echo.Context) error {
	treeName, err := url.QueryUnescape(c.Param("tree"))
	if err != nil {
		return err
	}

	treeFileName := fmt.Sprintf("attack_trees/descriptions/%s", treeName)
	treeFile, err := fs.GetFile(treeFileName)
	if err != nil {
		return err
	}

	treeContent, err := os.ReadFile(treeFile)
	if err != nil {
		return err
	}

	atkDir, err := fs.GetFile("attack_trees/descriptions/")
	if err != nil {
		return err
	}
	treeFiles, err := os.ReadDir(atkDir)
	if err != nil {
		return err
	}

	treeList := util.Map(treeFiles, func(t iofs.DirEntry) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: t.Name(),
			Link: fmt.Sprintf("/trees/%s", url.QueryEscape(t.Name())),
		}
	})

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape(treeFileName))

	var tree templates.TreeNode
	yaml.Unmarshal(treeContent, &tree)

	jsonTree, err := json.Marshal(&tree)
	if err != nil {
		fmt.Printf("LLL %s\n", err)
		return err
	}
	jsonStr := string(jsonTree)
	fmt.Println(jsonStr)

	return templates.Page(
		"Trees",
		"tree-editor", "Visual",
		func() templ.Component {
			return templates.FileList("/trees/", "attack_trees/descriptions", treeList)
		},
		func() templ.Component {
			return templates.TreeEditor("yaml", string(treeContent), saveEndpoint, &jsonStr)
		},
		func() templ.Component {
			return templates.Tree(url.QueryEscape(treeName), &tree)
		},
	).Render(c.Request().Context(), c.Response())
}

func EditTreeNode(c echo.Context) error {
	treeName, err := url.QueryUnescape(c.Param("tree"))
	if err != nil {
		return err
	}

	treeFileName := fmt.Sprintf("attack_trees/descriptions/%s", treeName)
	treeFile, err := fs.GetFile(treeFileName)
	if err != nil {
		return err
	}

	treeContent, err := os.ReadFile(treeFile)
	if err != nil {
		return err
	}

	atkDir, err := fs.GetFile("attack_trees/descriptions/")
	if err != nil {
		return err
	}
	treeFiles, err := os.ReadDir(atkDir)
	if err != nil {
		return err
	}

	treeList := util.Map(treeFiles, func(t iofs.DirEntry) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: t.Name(),
			Link: fmt.Sprintf("/trees/%s", url.QueryEscape(t.Name())),
		}
	})

	nodeFileName, err := url.QueryUnescape(c.Param("node"))
	if err != nil {
		return err
	}

	nodeFile, err := fs.GetFile(nodeFileName)
	if err != nil {
		return err
	}

	nodeContent, err := os.ReadFile(nodeFile)
	if err != nil {
		return err
	}

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape(nodeFileName))

	var tree templates.TreeNode
	yaml.Unmarshal(treeContent, &tree)

	newDesc := c.FormValue("description")
	if newDesc != "" {
		fmt.Printf("Changing description of '%s' to '%s', wrting to '%s'\n", nodeFileName, newDesc, treeFile)
		fs.ChangeTreeDescription(&tree, nodeFileName, newDesc)
		err := fs.SaveTreeDescription(&tree, treeFile)
		if err != nil {
			fmt.Println("ERROR")
			fmt.Println(err)
			return err
		}
	}

	return templates.Page(
		"Trees",
		"", "",
		func() templ.Component {
			return templates.SideBarList(treeList)
		},
		func() templ.Component { return templates.EditorComponent("yaml", string(nodeContent), saveEndpoint) },
		func() templ.Component {
			return templates.VerticalList(
				func() templ.Component { return templates.Tree(url.QueryEscape(treeName), &tree) },
				func() templ.Component {
					return templates.SideBarForm("#",
						templates.SideBarFormElement{
							Type:  "text",
							Id:    "description",
							Label: "Description",
						},
					)
				},
			)
		},
	).Render(c.Request().Context(), c.Response())
}

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
		err = os.WriteFile(urisFile, newContent, 0666)
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
		err = os.WriteFile(urisFile, newContent, 0666)
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
		"graphContainer", "Visualizer",
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

func ReasonerMainPage(c echo.Context) error {
	reasonDir, err := fs.GetFile("reasoner")
	if err != nil {
		return err
	}
	ruleFiles, err := os.ReadDir(reasonDir)
	if err != nil {
		return err
	}

	ruleList := util.Map(ruleFiles, func(t iofs.DirEntry) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: t.Name(),
			Link: fmt.Sprintf("/reasoner/%s", url.QueryEscape(t.Name())),
		}
	})

	return templates.Page(
		"Reasoner",
		"", "",
		func() templ.Component { return templates.FileList("/reasoner/", "reasoner", ruleList) },
		nil,
		nil,
	).Render(c.Request().Context(), c.Response())
}

func ReasonerRuleEditor(c echo.Context) error {
	ruleName, err := url.QueryUnescape(c.Param("rule"))
	if err != nil {
		return err
	}

	ruleFileName := fmt.Sprintf("reasoner/%s", ruleName)
	ruleFile, err := fs.GetFile(ruleFileName)
	if err != nil {
		return err
	}

	ruleContent, err := os.ReadFile(ruleFile)
	if err != nil {
		return err
	}

	reasonDir, err := fs.GetFile("reasoner")
	if err != nil {
		return err
	}
	ruleFiles, err := os.ReadDir(reasonDir)
	if err != nil {
		return err
	}

	ruleList := util.Map(ruleFiles, func(t iofs.DirEntry) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: t.Name(),
			Link: fmt.Sprintf("/reasoner/%s", url.QueryEscape(t.Name())),
		}
	})

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape(ruleFileName))

	return templates.Page(
		"Reasoner",
		"", "",
		func() templ.Component { return templates.FileList("/reasoner/", "reasoner/", ruleList) },
		func() templ.Component { return templates.EditorComponent("sparql", string(ruleContent), saveEndpoint) },
		nil,
	).Render(c.Request().Context(), c.Response())
}

func RegulationsMainPage(c echo.Context) error {
	regulationDirs, err := fs.GetRegulations()
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

	regName := c.Param("reg")

	regCfgFileName := fmt.Sprintf("regulations/%s/policies.yml", regName)
	regulationFile, err := fs.GetFile(regCfgFileName)
	if err != nil {
		return err
	}

	cfgContent, err := os.ReadFile(regulationFile)
	if err != nil {
		return err
	}

	regulationDirs, err := fs.GetRegulations()
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

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape(regCfgFileName))

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
	polName, err := url.QueryUnescape(c.Param("pol"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	regulationName := strings.Split(polName, "/")[1]
	regCfgFileName := fmt.Sprintf("regulations/%s/policies.yml", regulationName)
	regCfgFile, err := fs.GetFile(regCfgFileName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	regCfgContent, err := os.ReadFile(regCfgFile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	polFile, err := fs.GetFile(polName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	polContent, err := os.ReadFile(polFile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	regulationDirs, err := fs.GetRegulations()
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

	if formTitle != "" || formDescription != "" || formConsistency != "" || formViolations != "" || formMapping != "" {
		fmt.Printf("-> '%s' '%s' '%s' '%s' '%s'\n",
			formTitle,
			formDescription,
			formConsistency, // TODO: can't turn off consistency once turned on
			formViolations,
			formMapping,
		)

		policyRaw := util.First(policies, func(pol interface{}) bool {
			p := pol.(map[string]interface{})
			return p["file"] == polName
		})
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

		regulationFile, err := fs.GetFile(fmt.Sprintf("regulations/%s/policies.yml", regulationName))
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("Writing to '%s'\n", regulationFile)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = os.WriteFile(regulationFile, data, 0666)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	// TODO: tests
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
					Type:  templates.TEXT,
					Id:    "title",
					Label: "Title",
				},
				templates.SideBarFormElement{
					Type:  templates.TEXT,
					Id:    "description",
					Label: "Description",
				},
				templates.SideBarFormElement{
					Type:  templates.CHECKBOX,
					Id:    "consistency",
					Label: "Is consistency",
				},
				templates.SideBarFormElement{
					Type:  templates.TEXT,
					Id:    "violations",
					Label: "Maximum violations",
				},
				templates.SideBarFormElement{
					Type:  templates.TEXT,
					Id:    "mapping",
					Label: "Mapping message",
				},
			)
		},
	).Render(c.Request().Context(), c.Response())
}

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

func RequirementsMainPage(c echo.Context) error {
	requirementsFile, err := fs.GetFile("requirements/requirements.yml")
	if err != nil {
		return err
	}

	requriementsContent, err := os.ReadFile(requirementsFile)
	if err != nil {
		return err
	}

	contentList := []interface{}{}
	err = yaml.Unmarshal(requriementsContent, &contentList)
	if err != nil {
		return err
	}
	fmt.Println(string(requriementsContent))

	/*
		useCases := []templates.UseCase {}
		err = yaml.Unmarshal(requriementsContent, &useCases)
		if err != nil {
			return err
		}
	*/
	useCases := util.Map(contentList, func(ucRaw interface{}) templates.UseCase {
		uc := ucRaw.(map[string]interface{})

		useCase := uc["use case"].(string)
		isMisuseCase := uc["is misuse case"].(bool)
		reqsRaw := uc["requirements"].([]interface{})
		requirements := util.Map(reqsRaw, func(reqRaw interface{}) templates.Requirement {
			req := reqRaw.(map[string]interface{})

			title := req["title"].(string)
			description := req["description"].(string)
			query := req["query"].(string)

			return templates.Requirement{
				Title:       title,
				Description: description,
				Query:       query,
			}
		})

		return templates.UseCase{
			UseCase:      useCase,
			IsMisuseCase: isMisuseCase,
			Requirements: requirements,
		}
	})
	/*
		for _, uc := range useCases {
			fmt.Printf("'%s' '%v' '%d'\n", uc.UseCase, uc.IsMisuseCase, len(uc.Requirements))
		}
	*/

	rawJsonUCs, err := json.Marshal(&useCases)
	jsonUCs := string(rawJsonUCs)

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape("requirements/requirements.yml"))
	return templates.Page(
		"Requirements",
		"uc-editor", "Visual",
		func() templ.Component {
			return templates.UCSideBar(&useCases)
		},
		func() templ.Component {
			// return templates.EditorComponent("yaml", string(requriementsContent), saveEndpoint)
			return templates.UseCaseEditor("yaml", string(requriementsContent), saveEndpoint, &jsonUCs)
		},
		nil,
	).Render(c.Request().Context(), c.Response())
}

func RequirementEdit(c echo.Context) error {
	reqName, err := url.QueryUnescape(c.Param("req"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	requirementFile, err := fs.GetFile(reqName)
	if err != nil {
		return err
	}

	requriementQuery, err := os.ReadFile(requirementFile)
	if err != nil {
		return err
	}

	requirementsFile, err := fs.GetFile("requirements/requirements.yml")
	if err != nil {
		return err
	}

	requriementsContent, err := os.ReadFile(requirementsFile)
	if err != nil {
		return err
	}

	contentList := []interface{}{}
	err = yaml.Unmarshal(requriementsContent, &contentList)
	if err != nil {
		return err
	}
	fmt.Println(string(requriementsContent))

	useCases := util.Map(contentList, func(ucRaw interface{}) templates.UseCase {
		uc := ucRaw.(map[string]interface{})

		useCase := uc["use case"].(string)
		fmt.Printf("-- '%s'\n", useCase)
		isMisuseCase := uc["is misuse case"].(bool)
		reqsRaw := uc["requirements"].([]interface{})
		requirements := util.Map(reqsRaw, func(reqRaw interface{}) templates.Requirement {
			req := reqRaw.(map[string]interface{})

			title := req["title"].(string)
			description := req["description"].(string)
			query := req["query"].(string)

			return templates.Requirement{
				Title:       title,
				Description: description,
				Query:       query,
			}
		})

		return templates.UseCase{
			UseCase:      useCase,
			IsMisuseCase: isMisuseCase,
			Requirements: requirements,
		}
	})
	for _, uc := range useCases {
		fmt.Printf("'%s' '%v' '%d'\n", uc.UseCase, uc.IsMisuseCase, len(uc.Requirements))
	}

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape(reqName))
	return templates.Page(
		"Requirements",
		"", "",
		func() templ.Component {
			return templates.UCSideBar(&useCases)
		},
		func() templ.Component {
			return templates.EditorComponent("sparql", string(requriementQuery), saveEndpoint)
		},
		nil,
	).Render(c.Request().Context(), c.Response())
}

func RequirementDetails(c echo.Context) error {
	requirementsFile, err := fs.GetFile("requirements/requirements.yml")
	if err != nil {
		return err
	}

	requriementsContent, err := os.ReadFile(requirementsFile)
	if err != nil {
		return err
	}

	contentList := []interface{}{}
	err = yaml.Unmarshal(requriementsContent, &contentList)
	if err != nil {
		return err
	}

	/*
		useCases := []templates.UseCase {}
		err = yaml.Unmarshal(requriementsContent, &useCases)
		if err != nil {
			return err
		}
	*/
	useCases := util.Map(contentList, func(ucRaw interface{}) templates.UseCase {
		uc := ucRaw.(map[string]interface{})

		useCase := uc["use case"].(string)
		fmt.Printf("-- '%s'\n", useCase)
		isMisuseCase := uc["is misuse case"].(bool)
		reqsRaw := uc["requirements"].([]interface{})
		requirements := util.Map(reqsRaw, func(reqRaw interface{}) templates.Requirement {
			req := reqRaw.(map[string]interface{})

			title := req["title"].(string)
			description := req["description"].(string)
			query := req["query"].(string)

			return templates.Requirement{
				Title:       title,
				Description: description,
				Query:       query,
			}
		})

		return templates.UseCase{
			UseCase:      useCase,
			IsMisuseCase: isMisuseCase,
			Requirements: requirements,
		}
	})
	for _, uc := range useCases {
		fmt.Printf("%s %v %d\n", uc.UseCase, uc.IsMisuseCase, len(uc.Requirements))
	}

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape("requirements/requirements.yml"))
	return templates.Page(
		"Requirements",
		"", "",
		func() templ.Component {
			return templates.UCSideBar(&useCases)
		},
		func() templ.Component {
			return templates.EditorComponent("yaml", string(requriementsContent), saveEndpoint)
		},
		nil,
	).Render(c.Request().Context(), c.Response())

	/*
		return templates.Page(
			"Requirements",
		"",
			func() templ.Component {
				return templates.UCSideBar(&[]templates.UseCase{
					{Title: "a", IsMisuseCase: false, Requirements: []templates.Requirement{
						{Title: "b", Description: "b", Query: "b"},
					}},
				})
			},
			func() templ.Component { return templates.EditorComponent("yaml", "a: 1", "#") },
			func() templ.Component {
				return templates.UCDetails(
					"#",
					templates.UseCase{
						Title:        "",
						IsMisuseCase: false,
						Requirements: []templates.Requirement{},
					},
					templates.Requirement{
						Title:       "a",
						Description: "a",
						Query:       "a",
					},
				)
			},
		).Render(c.Request().Context(), c.Response())
	*/
}

func SchemasMainPage(c echo.Context) error {
	schemasDir, err := fs.GetFile("schemas")
	if err != nil {
		return err
	}

	schemaFiles, err := os.ReadDir(schemasDir)
	if err != nil {
		return err
	}

	schemas := util.Map(schemaFiles, func(s iofs.DirEntry) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: s.Name(),
			Link: fmt.Sprintf("/schemas/%s", s.Name()),
		}
	})

	return templates.Page(
		"Schemas",
		"", "",
		func() templ.Component { return templates.FileList("/schemas/", "schemas", schemas) },
		nil,
		nil,
	).Render(c.Request().Context(), c.Response())
}

func SchemaEditPage(c echo.Context) error {
	schemaName, err := url.QueryUnescape(c.Param("schema"))
	if err != nil {
		return err
	}

	schemaFile, err := fs.GetFile(fmt.Sprintf("schemas/%s", schemaName))
	if err != nil {
		return err
	}

	schemaContent, err := os.ReadFile(schemaFile)
	if err != nil {
		return err
	}

	schemasDir, err := fs.GetFile("schemas")
	if err != nil {
		return err
	}

	schemaFiles, err := os.ReadDir(schemasDir)
	if err != nil {
		return err
	}

	schemas := util.Map(schemaFiles, func(s iofs.DirEntry) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: s.Name(),
			Link: fmt.Sprintf("/schemas/%s", s.Name()),
		}
	})

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape(schemaFile))
	return templates.Page(
		"Schemas",
		"schemaEditorContainer", "Schema Editor",
		func() templ.Component { return templates.FileList("/schemas/", "schemas", schemas) },
		func() templ.Component {
			return templates.SchemaEditor("yaml", string(schemaContent), saveEndpoint)
		},
		nil,
	).Render(c.Request().Context(), c.Response())
}

func SaveEndpoint(c echo.Context) error {
	fmt.Println("SAVING")
	out, err := exec.Command("git", "-h").Output()
	if err != nil {
		fmt.Printf(">%s\n>%s\n", string(out), err)
		return err
	}

	userCookie := util.Filter(c.Request().Cookies(), func(cookie *http.Cookie) bool {
		return cookie.Name == "username"
	})[0]
	emailCookie := util.Filter(c.Request().Cookies(), func(cookie *http.Cookie) bool {
		return cookie.Name == "email"
	})[0]

	userName := userCookie.Value
	email := emailCookie.Value

	content, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	fName, err := url.QueryUnescape(c.Param("file"))
	if err != nil {
		return err
	}

	file, err := fs.GetFile(fName)
	if err != nil {
		return err
	}

	fmt.Printf("Writing to %s: %s \n", file, content)

	if err := os.WriteFile(file, content, 0666); err != nil {
		return err
	}

	fmt.Printf("In %s: %s %s\n", fs.LocalDir, userName, email)

	res, err := fs.AddAll(fs.LocalDir)
	if err != nil {
		fmt.Println(res)
		fmt.Println(err)
		return err
	}
	fmt.Sprintln(res)
	res, err = fs.Commit(fs.LocalDir, strconv.Itoa(rand.Int()))
	if err != nil {
		fmt.Println(res)
		fmt.Println(err)
		return err
	}
	fmt.Sprintln(res)

	return nil
}

func Analyse(c echo.Context) error {
	res, err := tool.Analyse("")

	fmt.Println(res)

	return err
}

func Test(c echo.Context) error {
	res, err := tool.Test()

	fmt.Println(res)

	return err
}

func DeleteFile(c echo.Context) error {
	file := c.QueryParam("path")
	path := fmt.Sprintf("'%s/%s'", fs.LocalDir, file)

	fmt.Printf("Delete '%s'\n", path)

	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func CreateFile(c echo.Context) error {
	file := c.QueryParam("path")
	path := fmt.Sprintf("'%s/%s'", fs.LocalDir, file)

	fmt.Printf("Create '%s'\n", path)

	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	return nil
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

	if err := os.WriteFile(fmt.Sprintf("%s/policies.yml", path), []byte("[]"), 0666); err != nil {
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
	userCookie := util.Filter(c.Request().Cookies(), func(cookie *http.Cookie) bool {
		return cookie.Name == "username"
	})[0]
	emailCookie := util.Filter(c.Request().Cookies(), func(cookie *http.Cookie) bool {
		return cookie.Name == "email"
	})[0]

	userName := userCookie.Value
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
	realRegPath, err := fs.GetFile(regPath)
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

	file, err := fs.GetFile(fName)
	if err != nil {
		return err
	}

	fmt.Printf("Writing to %s: %s \n", file, string(data))

	if err := os.WriteFile(file, data, 0666); err != nil {
		return err
	}

	fmt.Printf("In %s: %s %s\n", fs.LocalDir, userName, email)

	return nil
}

func UpdateTree(c echo.Context) error {
	userCookie := util.Filter(c.Request().Cookies(), func(cookie *http.Cookie) bool {
		return cookie.Name == "username"
	})[0]
	emailCookie := util.Filter(c.Request().Cookies(), func(cookie *http.Cookie) bool {
		return cookie.Name == "email"
	})[0]

	userName := userCookie.Value
	email := emailCookie.Value

	fName, err := url.QueryUnescape(c.Param("tree"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("Save to: %s\n", fName)

	body, err := io.ReadAll(c.Request().Body)

	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))

	var contents templates.TreeNode // interface{}
	err = json.Unmarshal(body, &contents)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// sync files with the config
	// file sync may be unnecessary, as once deleted files may be used in other queries
	// by the same assumptions that put all queries in the same directory, accessible by all trees
	/*
		pathComponents := strings.Split(fName, "/")
		regPath := strings.Join(pathComponents[:len(pathComponents)-1], "/") + "/consistency"
		realRegPath, err := fs.GetFile(regPath)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Printf("Syncing files at %s\n", realRegPath)
		err = util.SyncFileList(realRegPath, contents.GetQueries())
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("Files synced")
	*/
	// However, we still need to ensure all mentioned files exist
	queries := util.Map(contents.GetQueries(), func(path string) string {
		parts := strings.Split(path, "/")
		return parts[len(parts)-1]
	})
	dir, err := fs.GetFile("attack_trees/queries/")
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("It's here: %s\n", dir)

	list, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return err
	}

	dirContents := util.Map(list, func(d os.DirEntry) string {
		fmt.Println(d.Name())
		return d.Name()
	})

	for _, e := range queries {
		if !util.Contains(dirContents, e) {
			fmt.Printf("CREATING %s%s\n", dir, e)

			if err = os.WriteFile(fmt.Sprintf("%s%s", dir, e), []byte{}, 0666); err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	data, err := yaml.Marshal(contents)
	// _, err = yaml.Marshal(contents)
	// fmt.Printf("%+v\n", contents)

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Data to write \\/")
	fmt.Println(string(data))

	file, err := fs.GetFile(fName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Writing to %s: %s \n", file, string(data))

	if err := os.WriteFile(file, data, 0666); err != nil {
		return err
	}

	fmt.Printf("In %s: %s %s\n", fs.LocalDir, userName, email)

	return nil
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

func UpdateRequirements(c echo.Context) error {
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
	reqPath, err := fs.GetFile("requirements")
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

func UpdateURI(c echo.Context) error {
	desc, err := url.QueryUnescape(c.Param("desc"))
	if err != nil {
		return err
	}
	fmt.Printf("Desc: %s\n", desc)

	/*
		urisFile, err := fs.GetFile("uris.yml")
		if err != nil {
			return err
		}
		urisContent, err := os.ReadFile(urisFile)
		if err != nil {
			return err
		}

		uris := []interface{}{}
		err = yaml.Unmarshal(urisContent, &uris)
		if err != nil {
			return err
		}
	*/

	/*
		uriList := util.Map(uris, func(uri interface{}) string {
			uriObj := uri.(map[string]interface{})
			return uriObj["abreviation"].(string)
		})

		uri := c.FormValue("uri")
	*/
	uri := c.FormValue("uri")
	fmt.Printf("Old %s\n", uri)

	newURIAbrev := c.FormValue("abreviation")
	newURI := c.FormValue("new-uri")
	//if newURIAbrev != "" && newURI != "" {
	fmt.Printf("%s[%s]\n", newURIAbrev, newURI)
	//}

	return nil
}
