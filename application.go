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
	"github.com/plentiform/plentiform/models"
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

//func (app *Application) Render(w http.ResponseWriter, r *http.Request, name string, data pongo2.Context) error {
//func (app *Application) Render(w http.ResponseWriter, r *http.Request, name string, data context.Context) error {
func (app *Application) Render(w http.ResponseWriter, r *http.Request, name string, vars map[string]interface{}) error {
	//t, _ := pongo2.FromFile("templates/" + name + ".html")
	t, err := template.ParseFiles("templates/"+name+".html", "templates/layouts/public.html", "templates/gopher.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	type PageData struct {
		CurrentUser  *models.User
		Flashes      []interface{}
		RecaptchaKey string
	}
	//vars := map[string]interface{}{}

	tdata := new(PageData)

	session, _ := app.GetSession(r)

	if session.Values["userId"] != nil {
		user, _ := repo.NewUsersRepository(app.db).FindById(session.Values["userId"].(int))
		//data["currentUser"] = user
		//data = context.WithValue(data, "currentUser", user)
		tdata.CurrentUser = user
		vars["currentUser"] = user
	}

	//data["flashes"] = session.Flashes()
	//data = context.WithValue(data, "flashes", session.Flashes())
	tdata.Flashes = session.Flashes()
	vars["flashes"] = session.Flashes()
	//data["recaptcha_site_key"] = os.Getenv("RECAPTCHA_SITE_KEY")
	//data = context.WithValue(data, "recaptcha_site_key", os.Getenv("RECAPTCHA_SITE_KEY"))
	tdata.RecaptchaKey = os.Getenv("RECAPTCHA_SITE_KEY")
	vars["recaptcha_site_key"] = os.Getenv("RECAPTCHA_SITE_KEY")
	session.Save(r, w)

	//return t.ExecuteWriter(data, w)
	//return t.Execute(w, data)

	//return t.Execute(w, data.Value("recaptcha_site_key"))
	//data = context.WithValue(data, "currentUser", "Jim")
	//return t.Execute(w, tdata)
	return t.Execute(w, vars)
	//return t.Execute(w, r.WithContext(data))
	//return t.Execute(w, r)

}

func (app *Application) GetSession(r *http.Request) (*sessions.Session, error) {
	return app.sessions.Get(r, "plentiform")
}
