package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/plentiform/plentiform/mailers"
	"github.com/plentiform/plentiform/models"
	repo "github.com/plentiform/plentiform/repositories"
	"github.com/tuvistavie/securerandom"
)

func (app *Application) UsersNewHandler(w http.ResponseWriter, r *http.Request) {
	vars := map[string]interface{}{}
	app.Render(w, r, "users/new", vars)
}

func (app *Application) UsersCreateHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.GetSession(r)

	token, _ := securerandom.UrlSafeBase64(10, true)
	newUser := &models.User{
		Name:                   r.PostFormValue("name"),
		Email:                  r.PostFormValue("email"),
		PasswordDigest:         r.PostFormValue("password"),
		EmailConfirmationToken: &token,
	}

	if !app.recaptchaClient.Verify(*r) {
		session.AddFlash("Invalid ReCaptcha")
		session.Save(r, w)
		vars := map[string]interface{}{}
		vars["user"] = newUser
		app.Render(w, r, "users/new", vars)
		return
	}

	user, err := repo.NewUsersRepository(app.db).Create(newUser)
	if err != nil {
		log.Println(err)
		session.AddFlash("Woah, something bad happened.")
		session.Save(r, w)
		vars := map[string]interface{}{}
		vars["user"] = newUser
		app.Render(w, r, "users/new", vars)
		return
	}

	mailers.SendEmailConfirmation(app.emailClient, app.hostName, user)

	session.Values["userId"] = user.Id
	session.AddFlash(fmt.Sprintf("Welcome, %s!", user.Name))
	session.Save(r, w)

	http.Redirect(w, r, "/", 302)
}
