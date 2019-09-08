package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/haisum/recaptcha"
	_ "github.com/lib/pq"
	"github.com/sendgrid/sendgrid-go"
)

func Create() *Application {
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
