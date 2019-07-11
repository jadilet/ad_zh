package handlers

import (
	"net/http"

	"github.com/jadilet/ad_zh/models"

	"github.com/jadilet/ad_zh/session"
)

// ViewData send data to html template
type ViewData struct {
	Error string
	Data  string
}

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := store.Get(r, "cookie-user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentUser := session.GetCurrentUser(sess)

	if auth := currentUser.Authenticated; auth {
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}

	tpl.ExecuteTemplate(w, "login.tmpl", "")
}

func forgotPageHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "forgot.tmpl", "")
}

func resetPageHandler(w http.ResponseWriter, r *http.Request) {
	viewData := ViewData{Error: "", Data: r.URL.Query().Get("token")}
	tpl.ExecuteTemplate(w, "reset.tmpl", viewData)
}

func registerPageHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := store.Get(r, "cookie-user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentUser := session.GetCurrentUser(sess)

	if auth := currentUser.Authenticated; auth {
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}

	tpl.ExecuteTemplate(w, "register.tmpl", "")
}

func profileEditPageHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := store.Get(r, "cookie-user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentUser := session.GetCurrentUser(sess)

	if auth := currentUser.Authenticated; !auth {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	profileData := models.ProfileData{Error: "", User: currentUser.User}
	tpl.ExecuteTemplate(w, "profile-edit.tmpl", profileData)
}

func profilePageHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := store.Get(r, "cookie-user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentUser := session.GetCurrentUser(sess)

	if auth := currentUser.Authenticated; !auth {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	tpl.ExecuteTemplate(w, "profile.tmpl", currentUser.User)
}
