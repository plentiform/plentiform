package app

import (
	"database/sql"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/haisum/recaptcha"
	_ "github.com/lib/pq"
	repo "github.com/plentiform/plentiform/repositories"
	"github.com/sendgrid/sendgrid-go"
)

var CssHash string
var JsHash string

type Application struct {
	db              *sql.DB
	sessions        *sessions.CookieStore
	emailClient     *sendgrid.Client
	recaptchaClient recaptcha.R
	hostName        string
}

func (app *Application) Render(w http.ResponseWriter, r *http.Request, name string, vars map[string]interface{}) error {

	t, err := template.ParseFiles("templates/"+name+".html", "templates/layouts/public.html", "templates/gopher.html", "templates/forms/_form.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	session, _ := app.GetSession(r)

	if session.Values["userId"] != nil {
		user, _ := repo.NewUsersRepository(app.db).FindById(session.Values["userId"].(int))
		vars["currentUser"] = user
	}

	vars["css_hash"] = CssHash
	vars["js_hash"] = JsHash
	vars["session"] = session
	vars["recaptcha_site_key"] = os.Getenv("RECAPTCHA_SITE_KEY")
	vars["flashes"] = session.Flashes()
	session.Save(r, w)

	return t.Execute(w, vars)

}

func (app *Application) GetSession(r *http.Request) (*sessions.Session, error) {
	return app.sessions.Get(r, "plentiform")
}
