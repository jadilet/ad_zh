package handlers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jadilet/ad_zh/models"
	"github.com/jadilet/ad_zh/session"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const oauthGoogleURLAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8000/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	//create ouathState cookies
	ouathState := generateStateOauthCookie(w)

	u := googleOauthConfig.AuthCodeURL(ouathState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func oauthGoogleCallback(w http.ResponseWriter, r *http.Request, env *models.Env) {
	oauthState, _ := r.Cookie("oauthstate")
	sess, err := store.Get(r, "cookie-user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.FormValue("state") != oauthState.Value {
		log.Println("Invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := getUserDataFromGoogle(r.FormValue("code"))
	var userGoogle models.User

	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	json.Unmarshal([]byte(data), &userGoogle)

	log.Println(userGoogle.Name)
	log.Println(userGoogle.Email)
	log.Println(userGoogle.Picture)

	user, err := models.FindUser(env.DB, userGoogle.Email)

	switch {
	case err == sql.ErrNoRows:
		//user not found
		// register
		_, err = env.DB.Exec("INSERT INTO users(email, full_name, image_url) VALUES(?, ?, ?)", userGoogle.Email, userGoogle.Name, userGoogle.Picture)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal server error.", 500)
			return
		}

		// Save session
		sess, err := store.Get(r, "cookie-user")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newUser, err := models.FindUser(env.DB, userGoogle.Email)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if err == sql.ErrNoRows {
			http.Redirect(w, r, "/register", http.StatusFound)
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
	case err != nil:
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	default:
		// user found
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
	}
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	// Use code to take token and get user info from google
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Code exchange wrong: %s", err.Error())
	}

	response, err := http.Get(oauthGoogleURLAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed getting user info from google: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed read response: %s", err.Error())
	}

	return contents, nil
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(60 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}
