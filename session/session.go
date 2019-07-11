package session

import (
	"github.com/gorilla/sessions"

	"github.com/jadilet/ad_zh/models"
)

// CurrentUser holds a users account information
type CurrentUser struct {
	User          models.User
	Authenticated bool
}

// GetCurrentUser retrieve user from session
func GetCurrentUser(s *sessions.Session) CurrentUser {
	var user = CurrentUser{}

	val := s.Values["user"]

	// TODO implement check current user is exist in DB

	user, ok := val.(CurrentUser)
	if !ok {
		return CurrentUser{Authenticated: false}
	}

	return user
}
