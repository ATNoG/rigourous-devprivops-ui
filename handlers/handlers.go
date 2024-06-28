package handlers

import (
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
		func() templ.Component { return templates.SideBarList([]templates.SideBarListElement{{"Link", "#"}}) },
		func() templ.Component { return templates.LoginForm() },
		func() templ.Component {
			return templates.SideBarForm("Test", []templates.SideBarFormElement{{templates.TEXT, "idd", "label"}})
		},
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
		func() templ.Component { return templates.SideBarList([]templates.SideBarListElement{{"Link", "#"}}) },
		func() templ.Component { return templates.LoginForm() },
		func() templ.Component {
			return templates.SideBarForm("Test", []templates.SideBarFormElement{{templates.TEXT, "idd", "label"}})
		},
	).Render(c.Request().Context(), c.Response())
}

func DemoPage(c echo.Context) error {
	return templates.DemoPage().Render(c.Request().Context(), c.Response())
}

func TreesMainPage(c echo.Context) error {
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

	return templates.Page(
		"Trees",
		func() templ.Component {
			return templates.SideBarList(treeList)
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

	return templates.Page(
		"Trees",
		func() templ.Component {
			return templates.SideBarList(treeList)
		},
		func() templ.Component { return templates.EditorComponent("yaml", string(treeContent), saveEndpoint) },
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

	return templates.Page(
		"Trees",
		func() templ.Component {
			return templates.SideBarList(treeList)
		},
		func() templ.Component { return templates.EditorComponent("sparql", string(nodeContent), saveEndpoint) },
		func() templ.Component {
			return templates.Tree(url.QueryEscape(treeName), &tree)
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
		func() templ.Component {
			return templates.SideBarList(descriptions)
		},
		func() templ.Component { return templates.EditorComponent("yaml", "a: 1", "#") },
		nil,
	).Render(c.Request().Context(), c.Response())
}

func DescriptionEdit(c echo.Context) error {
	cookie := util.Filter(c.Request().Cookies(), func(cookie *http.Cookie) bool {
		return cookie.Name == "username"
	})[0]

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

	descriptions := util.Map(descs, func(d string) templates.SideBarListElement {
		return templates.SideBarListElement{
			Text: d,
			Link: fmt.Sprintf("/descriptions/%s", url.QueryEscape(d)),
		}
	})

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape(desc))

	return templates.Page(
		"My page",
		func() templ.Component {
			return templates.SideBarList(descriptions)
		},
		func() templ.Component { return templates.EditorComponent("yaml", string(descContent), saveEndpoint) },
		nil,
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
		func() templ.Component { return templates.SideBarList(ruleList) },
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
		func() templ.Component { return templates.SideBarList(ruleList) },
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
		func() templ.Component {
			return templates.SideBarList(regulations)
		},
		nil,
		nil,
	).Render(c.Request().Context(), c.Response())
}

func RegulationView(c echo.Context) error {

	regName := c.Param("reg")

	regCfgFileName := fmt.Sprintf("regulations/%s/policies.yml", regName)
	treeFile, err := fs.GetFile(regCfgFileName)
	if err != nil {
		return err
	}

	cfgContent, err := os.ReadFile(treeFile)
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

	policyFiles := util.Map(policies, func(pol interface{}) templates.SideBarListElement {
		p := pol.(map[string]interface{})
		fmt.Printf("!! /policy/%s\n", url.QueryEscape(p["file"].(string)))
		return templates.SideBarListElement{
			Text: p["file"].(string),
			Link: fmt.Sprintf("/policy/%s", url.QueryEscape(p["file"].(string))),
		}
	})

	regulations = append(regulations, policyFiles...)

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape(regCfgFileName))

	// TODO: tests
	return templates.Page(
		"Regulations",
		func() templ.Component {
			return templates.SideBarList(regulations)
		},
		func() templ.Component { return templates.EditorComponent("yaml", string(cfgContent), saveEndpoint) },
		func() templ.Component {
			return templates.SideBarForm("#", []templates.SideBarFormElement{
				{
					Type:  templates.TEXT,
					Id:    "Title",
					Label: "title",
				},
				{
					Type:  templates.TEXT,
					Id:    "Description",
					Label: "description",
				},
				{
					Type:  templates.CHECKBOX,
					Id:    "Is consistency",
					Label: "consistency",
				},
				{
					Type:  templates.TEXT,
					Id:    "Maximum violations",
					Label: "violations",
				},
				{
					Type:  templates.TEXT,
					Id:    "Mapping message",
					Label: "mapping",
				},
			})
		},
	).Render(c.Request().Context(), c.Response())
}

func PolicyEdit(c echo.Context) error {
	polName, err := url.QueryUnescape(c.Param("pol"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	regCfgFileName := fmt.Sprintf("regulations/%s/policies.yml", strings.Split(polName, "/")[1])
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

	// TODO: tests
	return templates.Page(
		"Regulations",
		func() templ.Component {
			return templates.SideBarList(regulations)
		},
		func() templ.Component { return templates.EditorComponent("sparql", string(polContent), saveEndpoint) },
		func() templ.Component {
			return templates.SideBarForm("#", []templates.SideBarFormElement{
				{
					Type:  templates.TEXT,
					Id:    "Title",
					Label: "title",
				},
				{
					Type:  templates.TEXT,
					Id:    "Description",
					Label: "description",
				},
				{
					Type:  templates.CHECKBOX,
					Id:    "Is consistency",
					Label: "consistency",
				},
				{
					Type:  templates.TEXT,
					Id:    "Maximum violations",
					Label: "violations",
				},
				{
					Type:  templates.TEXT,
					Id:    "Mapping message",
					Label: "mapping",
				},
			})
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
		func() templ.Component {
			return templates.SideBarList(extraDataList)
		},
		func() templ.Component {
			return templates.EditorComponent("yaml", string(extraDataContent), saveEndpoint)
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
		func() templ.Component {
			return templates.SideBarList(extraDataList)
		},
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

	var useCases []templates.UseCase
	err = yaml.Unmarshal(requriementsContent, &useCases)
	if err != nil {
		return err
	}
	for _, uc := range useCases {
		fmt.Printf("%s %v %d\n", uc.Title, uc.IsMisuseCase, len(uc.Requirements))
	}

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape("requirements/requirements.yml"))
	return templates.Page(
		"Requirements",
		func() templ.Component {
			return templates.UCSideBar(&useCases)
		},
		func() templ.Component {
			return templates.EditorComponent("yaml", string(requriementsContent), saveEndpoint)
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

	var useCases []templates.UseCase
	err = yaml.Unmarshal(requriementsContent, &useCases)
	if err != nil {
		return err
	}
	for _, uc := range useCases {
		fmt.Printf("%s %v %d\n", uc.Title, uc.IsMisuseCase, len(uc.Requirements))
	}

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape("requirements/requirements.yml"))
	return templates.Page(
		"Requirements",
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
		func() templ.Component {
			return templates.SideBarList(schemas)
		},
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
		func() templ.Component {
			return templates.SideBarList(schemas)
		},
		func() templ.Component {
			return templates.EditorComponent("yaml", string(schemaContent), saveEndpoint)
		},
		nil,
	).Render(c.Request().Context(), c.Response())
}

func SaveEndpoint(c echo.Context) error {
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

	desc, err := url.QueryUnescape(c.Param("file"))
	if err != nil {
		return err
	}

	descFile, err := fs.GetFile(desc)
	if err != nil {
		return err
	}

	fmt.Printf("Writing to %s: %s \n", descFile, content)

	if err := os.WriteFile(descFile, content, 0666); err != nil {
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
