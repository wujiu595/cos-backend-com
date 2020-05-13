package session

import (
	"github.com/wujiu2020/strip/sessions"
)

func SessionManager(provider sessions.SessionProvider) *sessions.SessionManager {
	manager := sessions.NewSessionManager(provider)
	return manager
}
