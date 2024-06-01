package main

import (
	"errors"
	"html/template"
	"microservice/internal/client"
	"net/http"
	"os"
)

func (app *application) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.clientError(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	url := os.Getenv("ENDPOINT_URL")
	if url == "" {
		app.serverError(w, errors.New("the endpoint url is undefined"))
		return
	}
	// GET endpoint data
	insurance, err := client.GetData(url)
	if err != nil {
		app.serverError(w, err)
	}

	ts, err := template.ParseFiles("./ui/html/homepage.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "homepage", insurance.InsuranceData)
	if err != nil {
		app.serverError(w, err)
	}
}
