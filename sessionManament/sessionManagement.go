package sessionmanament

import (
	"fmt"
	"sync"
)

type SessionManager struct {
	m                sync.Mutex
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
	res, ok := sm.sessionBranchMap[sessionKey]
	sm.m.Unlock()

	return res, ok
}

func (sm *SessionManager) AddSession(sessionKey string, branch string) {
	sm.m.Lock()
	fmt.Printf("%s:%s\n", sessionKey, branch)
	sm.sessionBranchMap[sessionKey] = branch
	sm.m.Unlock()
}
