package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
	"github.com/Joao-Felisberto/devprivops-ui/util"
	"github.com/labstack/echo"
)

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

	if err := fs.WriteFileSync(file, content, 0666); err != nil {
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
