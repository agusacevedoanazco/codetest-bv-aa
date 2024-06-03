package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func mockEndpoint(t *testing.T) (string, func()) {
	t.Helper()

	data := `{
		"insurance": {
			"name": "Test",
			"description": "Test description",
			"price": "100",
			"image": "image_url"
		}
	}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, data)
	}))

	return ts.URL, func() {
		ts.Close()
	}
}

func TestMain(t *testing.T) {
	testCases := []struct {
		name    string
		method  string
		path    string
		expCode int
	}{
		{name: "GetRoot", method: http.MethodGet, path: "/", expCode: http.StatusOK},
		{name: "PostRoot", method: http.MethodPost, path: "/", expCode: http.StatusMethodNotAllowed},
		{name: "GetStatic", method: http.MethodGet, path: "/static/css/main.css", expCode: http.StatusOK},
		{name: "NotFound", method: http.MethodGet, path: "/home", expCode: http.StatusNotFound},
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, tc.path, nil)

			url, cleanup := mockEndpoint(t)
			defer cleanup()

			os.Setenv("ENDPOINT_URL", url)
			defer os.Unsetenv("ENDPOINT_URL")

			handler := app.routes()

			handler.ServeHTTP(rec, req)

			if rec.Code != tc.expCode {
				t.Errorf("expected status %d, got %d", tc.expCode, rec.Code)
			}
		})
	}
}
