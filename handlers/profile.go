package handlers

import (
	"database/sql"
	"html"
	"log"
	"net/http"

	"github.com/jadilet/ad_zh/models"
	"github.com/jadilet/ad_zh/session"
)

func editProfile(w http.ResponseWriter, r *http.Request, env *models.Env) {
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

	r.ParseForm()

	email := html.EscapeString(r.FormValue("email"))
	address := html.EscapeString(r.FormValue("address"))
	fullName := html.EscapeString(r.FormValue("full_name"))
	telephone := html.EscapeString(r.FormValue("telephone"))

	var user models.User
	var errFind error

	// Find user by email
	user, errFind = models.FindUser(env.DB, currentUser.User.Email)

	switch {
	case errFind == sql.ErrNoRows:
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	case errFind != nil:
		log.Println(errFind.Error())
		http.Error(w, "Internal server error.", 500)
		return
	default:
		// Update Profile
		updateUser := models.User{Email: email, Name: fullName, Address: address, Telephone: telephone}
		err = models.UpdateUser(env.DB, updateUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err = models.FindUser(env.DB, updateUser.Email)
		switch {
		case err == sql.ErrNoRows:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		case err != nil:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		default:
			// Update session
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

			profileData := models.ProfileData{Error: "Successfully profile updated", User: user}
			tpl.ExecuteTemplate(w, "profile-edit.tmpl", profileData)
		}
	}
}
