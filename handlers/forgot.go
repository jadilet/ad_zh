package handlers

import (
	"database/sql"
	"fmt"
	"html"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/dchest/passwordreset"
	"github.com/jadilet/ad_zh/models"
	"github.com/jordan-wright/email"
	"golang.org/x/crypto/bcrypt"
)

// TODO set in enviroment variable
var secret = "lt85s1S6D9RhyA85jP6EakDb9b9j8Ics"

func reset(w http.ResponseWriter, r *http.Request, env *models.Env) {
	r.ParseForm()
	token := html.EscapeString(r.FormValue("token"))
	user, err := models.FindUserByToken(env.DB, token)

	switch {
	case err == sql.ErrNoRows:
		tpl.ExecuteTemplate(w, "forgot.tmpl", "Invalid token")
		return
	case err != nil:
		tpl.ExecuteTemplate(w, "forgot.tmpl", "Invalid token")
		return
	default:
		getPasswordHash := func(s string) ([]byte, error) {
			return []byte(user.Password), nil
		}
		_, err := passwordreset.VerifyToken(token, getPasswordHash, []byte(secret))
		if err != nil {
			// verification failed, don't allow password reset
			tpl.ExecuteTemplate(w, "forgot.tmpl", "Token verifivation failed")
			return
		}
		// OK, reset password for login (e.g. allow to change it)
		newPassword := html.EscapeString(r.FormValue("new_password"))
		confirmPassword := html.EscapeString(r.FormValue("confirm_password"))

		if newPassword != confirmPassword {
			viewData := ViewData{Error: "Mismatch new password and confirm password", Data: token}
			tpl.ExecuteTemplate(w, "reset.tmpl", viewData)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error, unable to generate hash password.", 500)
			return
		}

		user.Password = string(hashedPassword)

		err = models.UpdatePasswordUser(env.DB, user)

		if err != nil {
			log.Println(err.Error())

		}

		tpl.ExecuteTemplate(w, "login.tmpl", "Successfully reset password")
	}

}

func forgot(w http.ResponseWriter, r *http.Request, env *models.Env) {
	r.ParseForm()

	email := html.EscapeString(r.FormValue("email"))

	user, err := models.FindUser(env.DB, email)

	switch {
	case err == sql.ErrNoRows:
		tpl.ExecuteTemplate(w, "forgot.tmpl", fmt.Sprintf("%s email not found!", email))
		return
	case err != nil:
		log.Println(err.Error())
		http.Error(w, "Server error, unable to check data exist.", 500)
		return
	default:
		// generate token
		// send token with email
		token := passwordreset.NewToken(email, 1*time.Hour, []byte(user.Password), []byte(secret))
		sendMessage("mailgun@muslimlineru.ru", email, token)
		user.Token = token
		// Save token to db
		// TODO check error
		models.SaveToken(env.DB, user)

		tpl.ExecuteTemplate(w, "forgot.tmpl", fmt.Sprintf("Password reset link sent to %s", email))
	}
}

func sendMessage(sender, recipient, token string) {
	e := email.NewEmail()
	e.From = sender
	e.To = []string{recipient}
	e.Subject = "Password reset AD_ZH"
	e.Text = []byte(fmt.Sprintf("Someone requested a password reset for your account. If this was not you, please disregard this email. If you'd like to continue click the link. https://thawing-ridge-90494.herokuapp.com/reset?token=%s", token))
	err := e.Send("smtp.mailgun.org:587", smtp.PlainAuth("", os.Getenv("MAILGUN_USERNAME"), os.Getenv("MAILGUN_PASSWORD"), "smtp.mailgun.org"))
	if err != nil {
		panic(err)
	}
}
