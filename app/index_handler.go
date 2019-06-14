package app

import (
	"net/http"
)

func (app *Application) IndexHandler(w http.ResponseWriter, r *http.Request) {
	vars := map[string]interface{}{}
	app.Render(w, r, "index", vars)
}
