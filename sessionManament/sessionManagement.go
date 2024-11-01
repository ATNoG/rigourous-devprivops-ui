package sessionmanament

// https://github.com/gorilla/sessions

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/markbates/goth/gothic"
)

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

func (sm *SessionManager) GetBranch(sessionKey string) (string, bool) {
	sm.m.Lock()

	/*
		session, _ := store.Get(r, "session-name")
		// Set some session values.
		session.Values["foo"] = "bar"
		session.Values[42] = 43
	*/

	res, ok := sm.sessionBranchMap[sessionKey]
	sm.m.Unlock()

	return res, ok
}

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

func GetUserAndEmail(c *echo.Context) (string, string, error) {
	token, ok := (*c).Get("user").(*jwt.Token) // by default token is stored under `user` key
	if !ok {
		return "", "", errors.New("JWT token missing or invalid")
	}
	claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
	if !ok {
		return "", "", errors.New("failed to cast claims as jwt.MapClaims")
	}

	if err := claims.Valid(); err != nil {
		return "", "", err
	}

	return claims["username"].(string), claims["email"].(string), nil
}
