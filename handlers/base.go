package handlers

import (
	"encoding/gob"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jadilet/ad_zh/models"
	"github.com/jadilet/ad_zh/session"
)

// tpl holds all parsed templates
var tpl *template.Template

// store will hold all session data
var store = sessions.NewCookieStore([]byte(os.Getenv("COOKIE_SECRET_KEY")))

// InitSession flashes
func InitSession() {

	gob.Register(session.CurrentUser{})
}

// New Outh handlers
func New(env *models.Env) http.Handler {

	mux := mux.NewRouter()
	tpl = template.Must(template.ParseGlob("templates/*.tmpl"))
	InitSession()

	// OauthGoogle
	mux.HandleFunc("/auth/google/login", oauthGoogleLogin)
	mux.HandleFunc("/auth/google/callback", func(w http.ResponseWriter, r *http.Request) {
		oauthGoogleCallback(w, r, env)
	})

	// View
	mux.HandleFunc("/", indexPageHandler).Methods("GET")
	mux.HandleFunc("/login", indexPageHandler).Methods("GET")
	mux.HandleFunc("/register", registerPageHandler).Methods("GET")
	mux.HandleFunc("/profile", profilePageHandler).Methods("GET")
	mux.HandleFunc("/profile/edit", profileEditPageHandler).Methods("GET")

	// Authentication
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		login(w, r, env)
	}).Methods("POST")
	mux.HandleFunc("/logout", logout)

	// Registration
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		register(w, r, env)
	}).Methods("POST")

	// Profile
	mux.HandleFunc("/profile/edit", func(w http.ResponseWriter, r *http.Request) {
		editProfile(w, r, env)
	}).Methods("POST")

	return mux
}
