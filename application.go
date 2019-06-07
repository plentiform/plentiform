package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/haisum/recaptcha"
	_ "github.com/lib/pq"
	repo "github.com/plentiform/plentiform/repositories"
	"github.com/sendgrid/sendgrid-go"
)

type Application struct {
	db              *sql.DB
	sessions        *sessions.CookieStore
	emailClient     *sendgrid.Client
	recaptchaClient recaptcha.R
	hostName        string
}

func NewApplication() *Application {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	sessions := sessions.NewCookieStore([]byte(os.Getenv("SECRET_TOKEN")))
	emailClient := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	recaptchaClient := recaptcha.R{Secret: os.Getenv("RECAPTCHA_SECRET_TOKEN")}
	hostName := os.Getenv("HOST")

	return &Application{
		db:              db,
		sessions:        sessions,
		emailClient:     emailClient,
		recaptchaClient: recaptchaClient,
		hostName:        hostName,
	}
}

func (app *Application) Render(w http.ResponseWriter, r *http.Request, name string, vars map[string]interface{}) error {
	t, err := template.ParseFiles("templates/"+name+".html", "templates/layouts/public.html", "templates/gopher.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	session, _ := app.GetSession(r)

	if session.Values["userId"] != nil {
		user, _ := repo.NewUsersRepository(app.db).FindById(session.Values["userId"].(int))
		vars["currentUser"] = user
	}

	vars["session"] = session
	vars["recaptcha_site_key"] = os.Getenv("RECAPTCHA_SITE_KEY")
	vars["flashes"] = session.Flashes()
	session.Save(r, w)

	return t.Execute(w, vars)

}

func (app *Application) GetSession(r *http.Request) (*sessions.Session, error) {
	return app.sessions.Get(r, "plentiform")
}
