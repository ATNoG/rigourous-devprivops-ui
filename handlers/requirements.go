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

	file, err := fs.GetFile("requirements/requirements.yml")
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
