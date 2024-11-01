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
	"github.com/Joao-Felisberto/devprivops-ui/objects"
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/Joao-Felisberto/devprivops-ui/util"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"
)

func RequirementsMainPage(c echo.Context) error {
	userCookie, err := c.Cookie("username")
	if err != nil {
		return templates.Redirect("/").Render(c.Request().Context(), c.Response())
	}
	userName := userCookie.Value

	requirementsFile, err := fs.GetFile("requirements/requirements.yml", userName)
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
	fmt.Println("raw file \\/")
	fmt.Println(string(requriementsContent))

	fmt.Println("unmarshaled \\/")
	fmt.Printf("%+v\n", contentList)

	useCases := util.Map(contentList, func(ucRaw interface{}) objects.UseCase {
		uc := ucRaw.(map[string]interface{})

		useCase := uc["use case"].(string)
		isMisuseCase := uc["is misuse case"].(bool)
		reqsRaw := uc["requirements"].([]interface{})
		requirements := util.Map(reqsRaw, func(reqRaw interface{}) objects.Requirement {
			req := reqRaw.(map[string]interface{})

			title := req["title"].(string)
			description := req["description"].(string)
			query := req["query"].(string)

			return objects.Requirement{
				Title:       title,
				Description: description,
				Query:       query,
			}
		})

		return objects.UseCase{
			UseCase:      useCase,
			IsMisuseCase: isMisuseCase,
			Requirements: requirements,
		}
	})

	rawJsonUCs, err := json.Marshal(&useCases)
	if err != nil {
		fmt.Println(err)
		return err
	}
	jsonUCs := string(rawJsonUCs)

	// saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape("requirements/requirements.yml"))
	saveEndpoint := "/save-requirements"
	return templates.Page(
		"Requirements",
		"uc-editor", "Visual",
		templates.REQUIREMENTS,
		func() templ.Component {
			return templates.UCSideBar(&useCases)
		},
		func() templ.Component {
			return templates.UseCaseEditor("yaml", string(requriementsContent), saveEndpoint, &jsonUCs)
		},
		nil,
	).Render(c.Request().Context(), c.Response())
}

func RequirementEdit(c echo.Context) error {
	userCookie, err := c.Cookie("username")
	if err != nil {
		return templates.Redirect("/").Render(c.Request().Context(), c.Response())
	}
	userName := userCookie.Value

	reqName, err := url.QueryUnescape(c.Param("req"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	// fmt.Printf(">>> %s\n", reqName)

	requirementFile, err := fs.GetFile(reqName, userName)
	if err != nil {
		return err
	}

	requriementQuery, err := os.ReadFile(requirementFile)
	if err != nil {
		return err
	}

	requirementsFile, err := fs.GetFile("requirements/requirements.yml", userName)
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
	fmt.Println("raw file \\/")
	fmt.Println(string(requriementsContent))

	fmt.Println("unmarshaled \\/")
	fmt.Printf("%+v\n", contentList)

	var presentRequirement *objects.Requirement = nil

	useCases := util.Map(contentList, func(ucRaw interface{}) objects.UseCase {
		uc := ucRaw.(map[string]interface{})

		useCase := uc["use case"].(string)
		fmt.Printf("-- '%s'\n", useCase)
		isMisuseCase := uc["is misuse case"].(bool)
		reqsRaw := uc["requirements"].([]interface{})
		requirements := util.Map(reqsRaw, func(reqRaw interface{}) objects.Requirement {
			req := reqRaw.(map[string]interface{})

			title := req["title"].(string)
			description := req["description"].(string)
			query := req["query"].(string)

			res := objects.Requirement{
				Title:       title,
				Description: description,
				Query:       query,
			}

			if res.Query == reqName {
				presentRequirement = &res
			}

			return res
		})

		return objects.UseCase{
			UseCase:      useCase,
			IsMisuseCase: isMisuseCase,
			Requirements: requirements,
		}
	})
	for _, uc := range useCases {
		fmt.Printf("'%s' '%v' '%d'\n", uc.UseCase, uc.IsMisuseCase, len(uc.Requirements))
	}

	formTitle := c.FormValue("title")
	formDescription := c.FormValue("description")
	formQuery := c.FormValue("query")

	if formTitle != "" || formDescription != "" || formQuery != "" {
		if formTitle != "" {
			presentRequirement.Title = formTitle
		}
		if formDescription != "" {
			presentRequirement.Description = formDescription
		}
		if formQuery != "" {
			presentRequirement.Query = formQuery
		}

		data, err := yaml.Marshal(useCases)
		if err != nil {
			return err
		}

		requirementsFile, err := fs.GetFile("requirements/requirements.yml", userName)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("Writing to '%s'\n", requirementsFile)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = fs.WriteFileSync(requirementsFile, data, 0666)
		if err != nil {
			fmt.Println(err)
			return err
		}

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
	}

	saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape(reqName))
	return templates.Page(
		"Requirements",
		"", "",
		templates.REQUIREMENTS,
		func() templ.Component {
			return templates.UCSideBar(&useCases)
		},
		func() templ.Component {
			return templates.EditorComponent("sparql", string(requriementQuery), saveEndpoint)
		},
		func() templ.Component {
			return templates.SideBarForm(
				fmt.Sprintf("/requirements/%s", c.Param("req")),
				templates.SideBarFormElement{
					Type:    templates.TEXT,
					Id:      "title",
					Label:   "Title",
					Default: presentRequirement.Title,
				},
				templates.SideBarFormElement{
					Type:    templates.TEXT,
					Id:      "description",
					Label:   "Description",
					Default: presentRequirement.Description,
				},
			)
		},
	).Render(c.Request().Context(), c.Response())
}

func UpdateRequirements(c echo.Context) error {
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
	reqPath, err := fs.GetFile("requirements", userName)
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
