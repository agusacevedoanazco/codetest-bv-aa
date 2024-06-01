package main

import (
	"errors"
	"fmt"
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

	fmt.Fprintf(w, "%+v", insurance)
}
