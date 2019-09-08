package app

import (
	"html/template"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	repo "github.com/plentiform/plentiform/repositories"
)

var CssHash string
var JsHash string

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
