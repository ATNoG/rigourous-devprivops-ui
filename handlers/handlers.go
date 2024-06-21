package handlers

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/Joao-Felisberto/devprivops-ui/util"
	"github.com/a-h/templ"
	"github.com/labstack/echo"
)

func HomePage(c echo.Context) error {
	// return templates.DemoPage().Render(c.Request().Context(), c.Response())
	return templates.Page(
		"My page",
		func() templ.Component { return templates.SideBarList([]templates.SideBarListElement{{"Link", "#"}}) },
		func() templ.Component { return templates.EditorComponent("yaml", "a: 1", "#") },
		func() templ.Component {
			return templates.SideBarForm("Test", []templates.SideBarFormElement{{templates.TEXT, "idd", "label"}})
		},
	).Render(c.Request().Context(), c.Response())
}

func DemoPage(c echo.Context) error {
	return templates.DemoPage().Render(c.Request().Context(), c.Response())
}

func TreesMainPage(c echo.Context) error {
	return templates.Page(
		"Trees",
		func() templ.Component {
			return templates.SideBarList([]templates.SideBarListElement{
				{
					Text: "Tree 1",
					Link: "/trees/sample",
				},
			})
		},
		func() templ.Component { return templates.EditorComponent("yaml", "a: 1", "#") },
		nil,
	).Render(c.Request().Context(), c.Response())
}

func TreeView(c echo.Context) error {
	return templates.Page(
		"Trees",
		func() templ.Component {
			return templates.SideBarList([]templates.SideBarListElement{
				{
					Text: "Tree 1",
					Link: "/trees/sample",
				},
			})
		},
		func() templ.Component { return templates.EditorComponent("yaml", "a: 1", "#") },
		func() templ.Component {
			return templates.Tree(
				&templates.TreeNode{
					Description: "Root",
					Query:       "q",
					Children: []templates.TreeNode{
						{
							Description: "C1",
							Query:       "q",
							Children: []templates.TreeNode{
								{
									Description: "C11",
									Query:       "q",
									Children:    []templates.TreeNode{},
								},
							},
						},
						{
							Description: "C2",
							Query:       "q",
							Children:    []templates.TreeNode{},
						},
					},
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
		func() templ.Component {
			return templates.SideBarList(descriptions)
		},
		func() templ.Component { return templates.EditorComponent("yaml", "a: 1", "#") },
		func() templ.Component {
			return templates.PolicySideBar(&map[string][]templates.Policy{
				"Reg 1": {{"a", "a", "a", true, 1, "a"}},
			})
		},
	).Render(c.Request().Context(), c.Response())
}

func DescriptionEdit(c echo.Context) error {
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
		func() templ.Component {
			return templates.PolicySideBar(&map[string][]templates.Policy{
				"Reg 1": {{"a", "a", "a", true, 1, "a"}},
			})
		},
	).Render(c.Request().Context(), c.Response())
}

func ReasonerMainPage(c echo.Context) error {
	return templates.Page(
		"Reasoner",
		func() templ.Component { return templates.SideBarList([]templates.SideBarListElement{{"Rule 1", "#"}}) },
		func() templ.Component { return templates.EditorComponent("yaml", "a: 1", "#") },
		nil,
	).Render(c.Request().Context(), c.Response())
}

func RegulationsMainPage(c echo.Context) error {
	return templates.Page(
		"Regulations",
		func() templ.Component {
			return templates.SideBarList([]templates.SideBarListElement{
				{
					Text: "Regulation 1",
					Link: "/regulations/sample",
				},
			})
		},
		func() templ.Component { return templates.EditorComponent("yaml", "a: 1", "#") },
		nil,
	).Render(c.Request().Context(), c.Response())
}

func RegulationView(c echo.Context) error {
	// TODO: tests
	return templates.Page(
		"Regulations",
		func() templ.Component {
			return templates.SideBarList([]templates.SideBarListElement{
				{
					Text: "Regulation 1",
					Link: "/regulations/sample",
				},
			})
		},
		func() templ.Component { return templates.EditorComponent("yaml", "a: 1", "#") },
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
	return templates.Page(
		"Extra Data",
		func() templ.Component {
			return templates.SideBarList([]templates.SideBarListElement{
				{
					Text: "Query 1",
					Link: "extra-data/sample",
				},
			})
		},
		func() templ.Component { return templates.EditorComponent("yaml", "a: 1", "#") },
		nil,
	).Render(c.Request().Context(), c.Response())
}

func ExtraDataQuery(c echo.Context) error {
	return templates.Page(
		"Extra Data",
		func() templ.Component { return templates.SideBarList([]templates.SideBarListElement{{"Query 1", "#"}}) },
		func() templ.Component { return templates.EditorComponent("yaml", "a: 1", "#") },
		func() templ.Component {
			return templates.SideBarForm("#", []templates.SideBarFormElement{
				{
					Type:  templates.TEXT,
					Id:    "Heading",
					Label: "heading",
				},
				{
					Type:  templates.TEXT,
					Id:    "Description",
					Label: "description",
				},
				{
					Type:  templates.CHECKBOX,
					Id:    "Location",
					Label: "location",
				},
			})
		},
	).Render(c.Request().Context(), c.Response())
}

func RequirementsMainPage(c echo.Context) error {
	// return templates.DemoPage().Render(c.Request().Context(), c.Response())
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
		nil,
	).Render(c.Request().Context(), c.Response())
}

func RequirementDetails(c echo.Context) error {
	// return templates.DemoPage().Render(c.Request().Context(), c.Response())
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
}

func SchemasMainPage(c echo.Context) error {
	// return templates.DemoPage().Render(c.Request().Context(), c.Response())
	return templates.Page(
		"Schemas",
		func() templ.Component {
			return templates.SideBarList([]templates.SideBarListElement{{"Schema 1", "#"}})
		},
		func() templ.Component { return templates.EditorComponent("yaml", "a: 1", "#") },
		nil,
	).Render(c.Request().Context(), c.Response())
}

func SaveEndpoint(c echo.Context) error {
	/*
		var content []byte
		_, err := c.Request().Body.Read(content)
		if err != nil {
			return err
		}
	*/
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

	return nil
}
