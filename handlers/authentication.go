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

func login(w http.ResponseWriter, r *http.Request, env *models.Env) {
	r.ParseForm()
	sess, err := store.Get(r, "cookie-user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email := html.EscapeString(r.FormValue("email"))
	password := html.EscapeString(r.FormValue("password"))

	user, err := models.FindUser(env.DB, email)

	switch {
	case err == sql.ErrNoRows:
		tpl.ExecuteTemplate(w, "login.tmpl", fmt.Sprintf("User %s not found", email))
		return
	case err != nil:
		log.Println(err.Error())
		http.Error(w, "Server error, unable to check data exist.", 500)
		return
	default:
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err == nil {
			currentUser := &session.CurrentUser{
				User:          user,
				Authenticated: true,
			}

			sess.Values["user"] = currentUser

			err = sess.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/profile", http.StatusFound)

		} else {
			tpl.ExecuteTemplate(w, "login.tmpl", "Incorrect password ")
		}
	}

}

func logout(w http.ResponseWriter, r *http.Request) {
	sess, err := store.Get(r, "cookie-user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentUser := session.GetCurrentUser(sess)

	if currentUser.Authenticated {
		sess.Values["user"] = session.CurrentUser{}
		sess.Options.MaxAge = -1
	}

	err = sess.Save(r, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
