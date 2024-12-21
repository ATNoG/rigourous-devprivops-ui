package sessionmanament

// https://github.com/gorilla/sessions

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/markbates/goth/gothic"
)

var JWTSecret = ""

// JWT token generator
//
// `userID`: The user recieving the token
//
// `jwtSecret`: The secret with which to sign the token
//
// returns: The signed token string and an error if any happened during signing
func GenerateJWT(userID string, email string, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   time.Now().Add(time.Hour * time.Duration(3600)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// Finds the username of the current session inside the JWT token
//
// `c`: The current session
//
// returns: the username and an error, if the token does not exist or is invalid
func GetUsernameFromSession(c echo.Context) (string, error) {
	// Retrieve the token from cookies
	cookie, err := c.Cookie("auth")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", echo.NewHTTPError(http.StatusUnauthorized, "Missing cookie")
		}
		return "", echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Parse and validate the JWT token
	tokenString := cookie.Value
	tk, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid signing method")
		}
		return []byte(JWTSecret), nil
	})

	// Check for parsing errors or invalid tokens
	if err != nil || !tk.Valid {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	// Validate the token's expiration time
	claims, ok := tk.Claims.(jwt.MapClaims)
	if !ok || claims["exp"] == nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
	}

	// Check token expiration
	expirationTime := int64(claims["exp"].(float64))
	if time.Now().Unix() > expirationTime {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Token has expired")
	}

	user, ok := claims["sub"].(string)
	if !ok {
		return user, echo.NewHTTPError(http.StatusUnauthorized, "Token does not have a properly formated 'sub' field")
	}
	return user, nil
}

func GetEmailFromSession(c echo.Context) (string, error) {
	// Retrieve the token from cookies
	cookie, err := c.Cookie("auth")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", echo.NewHTTPError(http.StatusUnauthorized, "Missing cookie")
		}
		return "", echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Parse and validate the JWT token
	tokenString := cookie.Value
	tk, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid signing method")
		}
		return []byte(JWTSecret), nil
	})

	// Check for parsing errors or invalid tokens
	if err != nil || !tk.Valid {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	// Validate the token's expiration time
	claims, ok := tk.Claims.(jwt.MapClaims)
	if !ok || claims["exp"] == nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
	}

	// Check token expiration
	expirationTime := int64(claims["exp"].(float64))
	if time.Now().Unix() > expirationTime {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Token has expired")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return email, echo.NewHTTPError(http.StatusUnauthorized, "Token does not have a properly formated 'email' field")
	}
	return email, nil
}

type SessionManager struct {
	m                sync.Mutex
	jwt_secret       string
	sessionBranchMap map[string]string
}

func GetSessionManager() *SessionManager {
	return &SessionManager{
		m:                sync.Mutex{},
		sessionBranchMap: map[string]string{"master": "master"},
	}
}

/*
func (sm *SessionManager) GetBranch(sessionKey string) (string, bool) {

		sm.m.Lock()
		session, _ := store.Get(r, "session-name")
		// Set some session values.
		session.Values["foo"] = "bar"
		session.Values[42] = 43

		res, ok := sm.sessionBranchMap[sessionKey]
		sm.m.Unlock()


	res := fmt.Sprintf("%s/%s", fs.LocalDir, sessionKey)
	_, err := os.Stat(res)
	if err != nil {
		return "", false
	}

	return res, true
}
*/

func (sm *SessionManager) AddSession(r *http.Request, sessionKey string, branch string) error {
	sm.m.Lock()

	session, err := gothic.Store.Get(r, sessionKey)
	if err != nil {
		sm.m.Unlock()
		return err
	}
	session.Values["branch"] = branch

	fmt.Printf("%s:%s\n", sessionKey, branch)
	sm.sessionBranchMap[sessionKey] = branch

	sm.m.Unlock()
	return nil
}
