package handlers

import (
	"fmt"
	"html/template"
	"time"

	// "html/template"
	"net/http"
	"os"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
	sessionmanament "github.com/Joao-Felisberto/devprivops-ui/sessionManament"
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

// var JWTSecret = ""

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

	_, gh_key_found := os.LookupEnv("GITHUB_KEY")
	_, gh_sec_found := os.LookupEnv("GITHUB_SECRET")

	if gh_key_found && gh_sec_found {
		cookie, err := c.Cookie("ghAuth")
		if err != nil {
			if err == http.ErrNoCookie {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing cookie")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		// Parse and validate the JWT token
		tokenString := cookie.Value
		tk, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(sessionmanament.JWTSecret), nil
		})

		// Check for parsing errors or invalid tokens
		if err != nil || !tk.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		// Validate the token's expiration time
		claims, ok := tk.Claims.(jwt.MapClaims)
		if !ok || claims["exp"] == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
		}

		// Check token expiration
		expirationTime := int64(claims["exp"].(float64))
		if time.Now().Unix() > expirationTime {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token has expired")
		}
	}

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

	token, err := sessionmanament.GenerateJWT(userName, email, sessionmanament.JWTSecret)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	/*
		cookie := &http.Cookie{
			Name:     "auth",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * time.Duration(3600)),
			HttpOnly: true,
			Secure:   false, // Set to true in production (HTTPS)
			SameSite: http.SameSiteStrictMode,
		}
	*/
	cookie := new(http.Cookie)
	cookie.Name = "auth"
	cookie.Value = token
	cookie.Expires = time.Now().Add(time.Hour * time.Duration(3600))
	cookie.HttpOnly = true
	cookie.Secure = false
	cookie.SameSite = http.SameSiteStrictMode
	// http.SetCookie(c.Response(), cookie)
	c.SetCookie(cookie)

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
		fmt.Println("user auth completed")
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(res, gothUser)
	} else {
		fmt.Println("Could not complete user auth")
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
	fmt.Println("CALLBACK!")
	res := c.Response().Writer
	req := c.Request()

	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		fmt.Fprintln(res, err)
		return err
	}

	fmt.Printf("DATA: %+v\n", user)

	claims := jwt.MapClaims{
		"sub": user.AccessToken,
		"exp": time.Now().Add(time.Hour * time.Duration(3600)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tkStr, err := token.SignedString([]byte(sessionmanament.JWTSecret))
	if err != nil {
		return err
	}
	cookie := new(http.Cookie)
	cookie.Name = "ghAuth"
	cookie.Value = tkStr
	cookie.Expires = time.Now().Add(time.Hour * time.Duration(3600))
	cookie.HttpOnly = true
	cookie.Secure = false
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Path = "/credentials"
	c.SetCookie(cookie)

	return templates.Redirect(
		fmt.Sprintf("/credentials?username=%s&email=%s", user.NickName, user.Email),
	).Render(c.Request().Context(), c.Response())
}

/*
func DemoPage(c echo.Context) error {
	return templates.DemoPage().Render(c.Request().Context(), c.Response())
}
*/

// Function to verify whether a user is logged in.
// It serves as a middleware to ensure at each endpoint that the user accessing has a valid session.
//
// See https://pkg.go.dev/github.com/labstack/echo/v4#MiddlewareFunc
func EnsureLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Retrieve the token from cookies
		cookie, err := c.Cookie("auth")
		if err != nil {
			if err == http.ErrNoCookie {
				return c.Redirect(http.StatusTemporaryRedirect, "/")
			}
			fmt.Println(err)
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		// Parse and validate the JWT token
		tokenString := cookie.Value
		tk, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(sessionmanament.JWTSecret), nil
		})

		// Check for parsing errors or invalid tokens
		if err != nil || !tk.Valid {
			fmt.Println(err)
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		// Validate the token's expiration time
		claims, ok := tk.Claims.(jwt.MapClaims)
		if !ok || claims["exp"] == nil {
			fmt.Println(err)
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		// Check token expiration
		expirationTime := int64(claims["exp"].(float64))
		if time.Now().Unix() > expirationTime {
			fmt.Println(err)
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		// Store user information in the context
		/*
			userID := claims["sub"].(string)
			c.Set("userID", userID)
		*/

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
