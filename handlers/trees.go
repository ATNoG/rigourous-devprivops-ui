package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	iofs "io/fs"
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

func TreesMainPage(c echo.Context) error {
	userCookie, err := c.Cookie("username")
	if err != nil {
		return templates.Redirect("/").Render(c.Request().Context(), c.Response())
	}
	userName := userCookie.Value

	atkDir, err := fs.GetFile("attack_trees/descriptions/", userName)
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
		templates.TREES,
		func() templ.Component {
			return templates.FileList("/trees/", "attack_trees/descriptions", treeList)
		},
		nil,
		nil,
	).Render(c.Request().Context(), c.Response())
}

func TreeView(c echo.Context) error {
	userCookie, err := c.Cookie("username")
	if err != nil {
		return templates.Redirect("/").Render(c.Request().Context(), c.Response())
	}
	userName := userCookie.Value

	treeName, err := url.QueryUnescape(c.Param("tree"))
	if err != nil {
		return err
	}

	treeFileName := fmt.Sprintf("attack_trees/descriptions/%s", treeName)
	treeFile, err := fs.GetFile(treeFileName, userName)
	if err != nil {
		return err
	}

	treeContent, err := os.ReadFile(treeFile)
	if err != nil {
		return err
	}

	atkDir, err := fs.GetFile("attack_trees/descriptions/", userName)
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

	// saveEndpoint := fmt.Sprintf("/save/%s", url.QueryEscape(treeFileName))
	saveEndpoint := fmt.Sprintf("/save-tree/%s", url.QueryEscape(treeFileName))

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
		templates.TREES,
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
	userCookie, err := c.Cookie("username")
	if err != nil {
		return templates.Redirect("/").Render(c.Request().Context(), c.Response())
	}
	userName := userCookie.Value

	treeName, err := url.QueryUnescape(c.Param("tree"))
	if err != nil {
		return err
	}

	treeFileName := fmt.Sprintf("attack_trees/descriptions/%s", treeName)
	treeFile, err := fs.GetFile(treeFileName, userName)
	if err != nil {
		return err
	}

	treeContent, err := os.ReadFile(treeFile)
	if err != nil {
		return err
	}

	atkDir, err := fs.GetFile("attack_trees/descriptions/", userName)
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

	nodeFile, err := fs.GetFile(nodeFileName, userName)
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

	nodeDescription := tree.GetNodeDescription(nodeFileName)
	if nodeDescription == "" {
		fmt.Printf("Could not find '%s' in tree\n", nodeFileName)
		return fmt.Errorf("could not find '%s' in tree", nodeFileName)
	}

	return templates.Page(
		"Trees",
		"", "",
		templates.TREES,
		/*
			func() templ.Component {
				return templates.SideBarList(treeList)
			},
		*/
		func() templ.Component {
			return templates.FileList("/trees/", "attack_trees/descriptions", treeList)
		},
		func() templ.Component { return templates.EditorComponent("yaml", string(nodeContent), saveEndpoint) },
		func() templ.Component {
			return templates.VerticalList(
				func() templ.Component { return templates.Tree(url.QueryEscape(treeName), &tree) },
				func() templ.Component {
					return templates.SideBarForm("#",
						templates.SideBarFormElement{
							Type:    "text",
							Id:      "description",
							Label:   "Description",
							Default: nodeDescription,
						},
					)
				},
			)
		},
	).Render(c.Request().Context(), c.Response())
}

func UpdateTree(c echo.Context) error {
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
	dir, err := fs.GetFile("attack_trees/queries/", userName)
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

			if err = fs.WriteFileSync(fmt.Sprintf("%s%s", dir, e), []byte{}, 0666); err != nil {
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

	file, err := fs.GetFile(fName, userName)
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
