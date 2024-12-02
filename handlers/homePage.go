package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

func HomePage(c echo.Context) error {
	_, gh_key_found := os.LookupEnv("GITHUB_KEY")
	_, gh_sec_found := os.LookupEnv("GITHUB_SECRET")

	if gh_key_found && gh_sec_found {
		return templates.Redirect("/auth?provider=github").Render(c.Request().Context(), c.Response())
	} else {
		prevUser := c.QueryParam("username")
		prevMail := c.QueryParam("email")

		return templates.LoginPage(prevUser, prevMail).Render(c.Request().Context(), c.Response())
	}
}

func GetCredentials(c echo.Context) error {
	prevUser := c.QueryParam("username")
	prevMail := c.QueryParam("email")

	/*
		return templates.Page(
			"Home page",
			"", "",
			-1,
			nil,
			nil,
			nil,
		).Render(c.Request().Context(), c.Response())
	*/
	return templates.LoginPage(prevUser, prevMail).Render(c.Request().Context(), c.Response())
}

func SimpleLogIn(c echo.Context) error {
	userNameCookie := new(http.Cookie)
	userNameCookie.Name = "username"
	userName := c.FormValue("username")
	userNameCookie.Value = userName
	userNameCookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(userNameCookie)

	emailCookie := new(http.Cookie)
	emailCookie.Name = "email"
	email := c.FormValue("email")
	emailCookie.Value = email
	emailCookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(emailCookie)

	fmt.Println("HERE!")
	fs.SessionManager.AddSession(c.Request(), userName, userName)
	fs.SetupRepo(userName, userName, email)

	return templates.Page(
		"Home page",
		"", "",
		-1,
		nil,
		nil,
		nil,
	).Render(c.Request().Context(), c.Response())
}

func Login(c echo.Context) error {
	res := c.Response().Writer
	req := c.Request()

	if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(res, gothUser)
	} else {
		gothic.BeginAuthHandler(res, req)
	}

	return nil
}

func Logout(c echo.Context) error {
	err := gothic.Logout(c.Response().Writer, c.Request())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func Callback(c echo.Context) error {
	res := c.Response().Writer
	req := c.Request()

	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		fmt.Fprintln(res, err)
		return err
	}
	/*
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(res, user)
	*/
	fmt.Printf("DATA: %+v\n", user)

	/*
		return templates.Page(
			"Home page",
			"", "",
			-1,
			nil,
			nil,
			nil,
		).Render(c.Request().Context(), c.Response())
	*/
	return templates.Redirect(
		fmt.Sprintf("/credentials?username=%s&email=%s", user.NickName, user.Email),
	).Render(c.Request().Context(), c.Response())
}

func DemoPage(c echo.Context) error {
	return templates.DemoPage().Render(c.Request().Context(), c.Response())
}

var userTemplate = `
<p><a href="/logout?provider={{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`
