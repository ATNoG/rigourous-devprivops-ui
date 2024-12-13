package handlers

import (
	"fmt"
	"html/template"

	// "html/template"
	"net/http"
	"os"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

// Website entry point that redirects to Github OAuth or requests the login information itself.
//
// `c`: The echo context
//
// returns: error if any internal function, like file reading, or template rendering fails.
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

// Endpoint to get the git credentials from the user.
//
// `c`: The echo context
//
// returns: error if any internal function, like file reading, or template rendering fails.
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

// Endpoint of the main page
//
// `c`: The echo context
//
// returns: error if any internal function, like file reading, or template rendering fails.
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

// Endpoint to redirect to github oauth authentication
//
// `c`: The echo context
//
// returns: error if any internal function, like file reading, or template rendering fails.
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

// Endpoint to logout the currently logged in user
//
// `c`: The echo context
//
// returns: error if any internal function, like file reading, or template rendering fails.
func Logout(c echo.Context) error {
	err := gothic.Logout(c.Response().Writer, c.Request())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Endpoint to which the github authenticator redirects to after authenticating
//
// `c`: The echo context
//
// returns: error if any internal function, like file reading, or template rendering fails.
func Callback(c echo.Context) error {
	res := c.Response().Writer
	req := c.Request()

	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		fmt.Fprintln(res, err)
		return err
	}

	fmt.Printf("DATA: %+v\n", user)
	c.Set("user", user)

	return templates.Redirect(
		fmt.Sprintf("/credentials?username=%s&email=%s", user.NickName, user.Email),
	).Render(c.Request().Context(), c.Response())
}

/*
func DemoPage(c echo.Context) error {
	return templates.DemoPage().Render(c.Request().Context(), c.Response())
}
*/

func EnsureLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check if the user is authenticated
		user := c.Get("user")
		if user == nil {
			return c.Redirect(http.StatusFound, "/")
			// return echo.NewHTTPError(http.StatusUnauthorized, "You must be logged in to view this page.")
		}
		emailCookie, err := c.Cookie("email")
		if err != nil {
			return c.Redirect(http.StatusFound, "/")
		}
		email := emailCookie.Value
		if email == "" {
			return c.Redirect(http.StatusFound, "/")
		}
		return next(c)
	}
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
