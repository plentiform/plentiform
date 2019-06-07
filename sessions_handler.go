package main

import (
	"net/http"

	repo "github.com/plentiform/plentiform/repositories"
)

func (app *Application) SessionsNewHandler(w http.ResponseWriter, r *http.Request) {
	//app.Render(w, r, "index", pongo2.Context{})
	//app.Render(w, r, "index", r.Context())
	vars := map[string]interface{}{}
	app.Render(w, r, "index", vars)
}

func (app *Application) SessionsCreateHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.GetSession(r)

	user_table_exists := repo.NewUsersRepository(app.db).UserTableExists()
	if app.db.Ping() != nil {
		session.AddFlash("Can't connect to the database.")
		session.Save(r, w)
		//app.Render(w, r, "index", pongo2.Context{"email": r.PostFormValue("email")})
		//ctx := r.Context()
		//ctx = context.WithValue(ctx, "email", r.PostFormValue("email"))
		//app.Render(w, r, "index", ctx)
		vars := map[string]interface{}{}
		vars["email"] = r.PostFormValue("email")
		app.Render(w, r, "index", vars)
		return
	} else if user_table_exists == false {
		session.AddFlash("There is no user table. Make sure migrations have been run.")
		session.Save(r, w)
		//app.Render(w, r, "index", pongo2.Context{"email": r.PostFormValue("email")})
		//ctx := r.Context()
		//ctx = context.WithValue(ctx, "email", r.PostFormValue("email"))
		//app.Render(w, r, "index", ctx)
		vars := map[string]interface{}{}
		vars["email"] = r.PostFormValue("email")
		//vars["flashes"] = session.Flashes()
		app.Render(w, r, "index", vars)
		return
	}

	valid_email, _ := repo.NewUsersRepository(app.db).FindByEmail(r.PostFormValue("email"))
	if valid_email == nil {
		session.AddFlash("There is no account for the email: " + r.PostFormValue("email"))
		session.Save(r, w)
		//app.Render(w, r, "index", pongo2.Context{"email": r.PostFormValue("email")})
		//ctx := r.Context()
		//ctx = context.WithValue(ctx, "email", r.PostFormValue("email"))
		//app.Render(w, r, "index", ctx)
		vars := map[string]interface{}{}
		vars["email"] = r.PostFormValue("email")
		//vars["flashes"] = session.Flashes()
		app.Render(w, r, "index", vars)
		return
	}

	user, _ := repo.NewUsersRepository(app.db).FindByEmailAndPassword(r.PostFormValue("email"), r.PostFormValue("password"))
	if user == nil {
		session.AddFlash("Invalid password.")
		session.Save(r, w)
		//app.Render(w, r, "index", pongo2.Context{"email": r.PostFormValue("email")})
		//ctx := r.Context()
		//ctx = context.WithValue(ctx, "email", r.PostFormValue("email"))
		//app.Render(w, r, "index", ctx)
		vars := map[string]interface{}{}
		vars["email"] = r.PostFormValue("email")
		//vars["flashes"] = session.Flashes()
		app.Render(w, r, "index", vars)
		return
	}

	session.Values["userId"] = user.Id
	session.AddFlash("Successfully logged in!")
	session.Save(r, w)

	http.Redirect(w, r, "/", 302)
}

func (app *Application) SessionsDestroyHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.GetSession(r)
	session.Values["userId"] = nil
	session.AddFlash("Successfully logged out!")
	session.Save(r, w)

	http.Redirect(w, r, "/", 302)
}
