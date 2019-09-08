package app

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/haisum/recaptcha"
	_ "github.com/lib/pq"
	"github.com/sendgrid/sendgrid-go"
)

type Application struct {
	db              *sql.DB
	sessions        *sessions.CookieStore
	emailClient     *sendgrid.Client
	recaptchaClient recaptcha.R
	hostName        string
}

func (app *Application) GetSession(r *http.Request) (*sessions.Session, error) {
	return app.sessions.Get(r, "plentiform")
}
