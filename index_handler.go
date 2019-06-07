package main

import (
	"net/http"
)

func (app *Application) IndexHandler(w http.ResponseWriter, r *http.Request) {
	//app.Render(w, r, "index", pongo2.Context{})
	//app.Render(w, r, "index", context.Background())
	//app.Render(w, r, "index", r.Context())
	vars := map[string]interface{}{}
	app.Render(w, r, "index", vars)
}
