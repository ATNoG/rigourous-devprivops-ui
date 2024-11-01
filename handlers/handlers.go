package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"os"
	"strconv"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/labstack/echo/v4"
)

func SaveEndpoint(c echo.Context) error {
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

	content, err := io.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fName, err := url.QueryUnescape(c.Param("file"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	file, err := fs.GetFile(fName, userName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Writing to %s: %s \n", file, content)

	if err := fs.WriteFileSync(file, content, 0666); err != nil {
		fmt.Println(err)
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

func DeleteFile(c echo.Context) error {
	userCookie, err := c.Cookie("username")
	if err != nil {
		return templates.Redirect("/").Render(c.Request().Context(), c.Response())
	}
	userName := userCookie.Value
	branch, exists := fs.SessionManager.GetBranch(userName)
	if !exists {
		return fmt.Errorf("user is not logged in")
	}

	file := c.QueryParam("path")
	path := fmt.Sprintf("%s/%s/%s", fs.LocalDir, branch, file)

	fmt.Printf("Delete '%s'\n", path)

	err = os.RemoveAll(path)
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

func UpdateURI(c echo.Context) error {
	desc, err := url.QueryUnescape(c.Param("desc"))
	if err != nil {
		return err
	}
	fmt.Printf("Desc: %s\n", desc)

	uri := c.FormValue("uri")
	fmt.Printf("Old %s\n", uri)

	newURIAbrev := c.FormValue("abreviation")
	newURI := c.FormValue("new-uri")
	//if newURIAbrev != "" && newURI != "" {
	fmt.Printf("%s[%s]\n", newURIAbrev, newURI)
	//}

	return nil
}
