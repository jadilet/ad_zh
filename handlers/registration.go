package handlers

import (
	"database/sql"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/jadilet/ad_zh/session"

	"github.com/jadilet/ad_zh/models"
	"golang.org/x/crypto/bcrypt"
)

func register(w http.ResponseWriter, r *http.Request, env *models.Env) {
	r.ParseForm()

	email := html.EscapeString(r.FormValue("email"))
	password := html.EscapeString(r.FormValue("password"))
	passwordConfirm := html.EscapeString(r.FormValue("confirm_password"))

	if password != passwordConfirm {
		tpl.ExecuteTemplate(w, "register.tmpl", "Mismatch password and confirm password")
		return
	}

	var user models.User
	var errFind error

	// Find user by email
	user, errFind = models.FindUser(env.DB, email)

	switch {
	case errFind == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error, unable to generate hash password.", 500)
			return
		}

		_, err = env.DB.Exec("INSERT INTO users(email, password) VALUES(?, ?)", email, hashedPassword)
		if err != nil {
			http.Error(w, "Server error, unable to insert your data.", 500)
			return
		}

		// Save session
		sess, err := store.Get(r, "cookie-user")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newUser, err := models.FindUser(env.DB, email)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		currentUser := &session.CurrentUser{
			User:          newUser,
			Authenticated: true,
		}

		sess.Values["user"] = currentUser
		err = sess.Save(r, w)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/profile/edit", http.StatusFound)
		return
	case errFind != nil:
		log.Println(errFind.Error())
		http.Error(w, "Server error, unable to check data exist.", 500)
		return
	default:
		tpl.ExecuteTemplate(w, "register.tmpl", fmt.Sprintf("User %s already exist", user.Email))
	}
}
