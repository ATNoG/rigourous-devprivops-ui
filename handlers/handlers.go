package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
	sessionmanament "github.com/Joao-Felisberto/devprivops-ui/sessionManament"
	"github.com/labstack/echo/v4"
)

// Endpoint to save a generic file sent in the request body
//
// `c`: The echo context
//
// returns: error if any internal function, like file reading, or template rendering fails.
func SaveEndpoint(c echo.Context) error {
	userName, err := sessionmanament.GetUsernameFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	email, err := sessionmanament.GetEmailFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

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

// Endpoint to delete any generic file
//
// `c`: The echo context
//
// # QUERY PARAMETERS
//
// `path`: The path of the file to delete
//
// returns: error if any internal function, like file reading, or template rendering fails.
func DeleteFile(c echo.Context) error {
	userName, err := sessionmanament.GetUsernameFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	branch, exists := fs.GetBranch(userName)
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

// Endpoint to create an empty file at a specific location
//
// `c`: The echo context
//
// # QUERY PARAMETERS
//
// `path`: The path of the file to create
//
// returns: error if any internal function, like file reading, or template rendering fails.
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

// Endpoint to update a description file's default URI
//
// `c`: The echo context
//
// returns: error if any internal function, like file reading, or template rendering fails.
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
