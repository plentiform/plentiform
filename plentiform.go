package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	pipeline "github.com/plentiform/go-asset-pipeline"
	"github.com/plentiform/plentiform/app"
)

func main() {

	app := app.Create()

	r := mux.NewRouter()
	r.HandleFunc("/", app.IndexHandler).Methods("GET")
	r.HandleFunc("/login", app.SessionsNewHandler).Methods("GET")
	r.HandleFunc("/login", app.SessionsCreateHandler).Methods("POST")
	r.HandleFunc("/logout", app.SessionsDestroyHandler).Methods("GET")
	r.HandleFunc("/signup", app.UsersNewHandler).Methods("GET")
	r.HandleFunc("/signup", app.UsersCreateHandler).Methods("POST")
	r.HandleFunc("/f/{uuid}", app.SubmissionsCreateHandler).Methods("POST")
	r.HandleFunc("/email_confirmation/new", app.EmailConfirmationsNewHandler).Methods("GET")
	r.HandleFunc("/email_confirmation", app.EmailConfirmationsCreateHandler).Methods("POST")
	r.HandleFunc("/email_confirmation", app.EmailConfirmationsShowHandler).Methods("GET")
	r.HandleFunc("/forms/{formId:[0-9]+}/submissions/{submissionId:[0-9]+}", app.RequireAuthentication(app.RequireEmailConfirmation(app.SubmissionsDestroyHandler))).Methods("DELETE")
	r.HandleFunc("/forms/{id:[0-9]+}/edit", app.RequireAuthentication(app.RequireEmailConfirmation(app.FormsEditHandler))).Methods("GET")
	r.HandleFunc("/forms/{id:[0-9]+}", app.RequireAuthentication(app.RequireEmailConfirmation(app.FormsShowHandler))).Methods("GET")
	r.HandleFunc("/forms/{id:[0-9]+}", app.RequireAuthentication(app.RequireEmailConfirmation(app.FormsDestroyHandler))).Methods("DELETE")
	r.HandleFunc("/forms/{id:[0-9]+}", app.RequireAuthentication(app.RequireEmailConfirmation(app.FormsUpdateHandler))).Methods("POST")
	r.HandleFunc("/forms/new", app.RequireAuthentication(app.RequireEmailConfirmation(app.FormsNewHandler))).Methods("GET")
	r.HandleFunc("/forms", app.RequireAuthentication(app.RequireEmailConfirmation(app.FormsIndexHandler))).Methods("GET")
	r.HandleFunc("/forms", app.RequireAuthentication(app.RequireEmailConfirmation(app.FormsCreateHandler))).Methods("POST")

	setup_pipeline(r)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "3000"
	}

	log.Println("Listening on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port,
		handlers.CompressHandler(
			handlers.HTTPMethodOverrideHandler(
				handlers.LoggingHandler(os.Stdout, r)))))

}

func setup_pipeline(r *mux.Router) {

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./assets/public/")))
	// Asset pipeline to concat, minify, and fingerprint css & js
	result, _ := pipeline.Compile(pipeline.Config{
		Files: []string{
			"assets/css/main.css",
			"assets/css/components/*",
			"assets/js/main.js",
		},
		Minify:    true,
		OutputDir: "assets/public/",
	})
	app.CssHash = result[pipeline.CSS].Hash
	app.JsHash = result[pipeline.JS].Hash

}
