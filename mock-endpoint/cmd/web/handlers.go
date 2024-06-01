package main

import (
	"encoding/json"
	"net/http"
)

// Data structure
type StructData struct {
	Data Data `json:"insurance"`
}

type Data struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Image       string `json:"image"`
}

// Data to send
var data = StructData{
	Data: Data{
		Name:        "Nombre de Seguro",
		Description: "Descripcion del Seguro",
		Price:       "Valor",
		Image:       "url_de_imagen",
	},
}

// Log Request Middleware
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

// Route multiplexer
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.root)

	return app.logRequest(mux)
}

// Root Handler
func (app *application) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		app.serverError(w, err)
		return
	}
}
