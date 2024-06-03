package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetDataOK(t *testing.T) {
	mockResponse := Insurance{
		InsuranceData: InsuranceData{
			Name:        "Test name",
			Description: "Test Insurance long description",
			Price:       "Some value",
			Image:       "https://someimage-url.local/myimage.png",
		},
	}

	body, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}))
	defer server.Close()

	data, err := GetData(server.URL)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	if data.InsuranceData.Name != mockResponse.InsuranceData.Name {
		t.Errorf("expected name %s, got %s", mockResponse.InsuranceData.Name, data.InsuranceData.Name)
	}
	if data.InsuranceData.Description != mockResponse.InsuranceData.Description {
		t.Errorf("expected description %s, got %s", mockResponse.InsuranceData.Description, data.InsuranceData.Description)
	}
	if data.InsuranceData.Price != mockResponse.InsuranceData.Price {
		t.Errorf("expected price %s, got %s", mockResponse.InsuranceData.Price, data.InsuranceData.Price)
	}
	if data.InsuranceData.Image != mockResponse.InsuranceData.Image {
		t.Errorf("expected image url %s, got %s", mockResponse.InsuranceData.Image, data.InsuranceData.Image)
	}
}

func TestGetDataServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := GetData(server.URL)
	if err == nil {
		t.Errorf("expected error, got none")
	}

	expErrMsg := "unexpected client response 500"
	if err.Error() != expErrMsg {
		t.Errorf("expected error message %s, got %s", expErrMsg, err)
	}
}

func TestGetDataBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("super convincing json"))
	}))
	defer server.Close()

	_, err := GetData(server.URL)
	if err == nil {
		t.Fatalf("expected an error, got none")
	}

	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected a JSON syntax error, got %v", err)
	}
}

func TestGetDataTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Minute + time.Second) // sleep a second over the timeout of the client package
	}))
	defer server.Close()

	_, err := GetData(server.URL)
	if err == nil {
		t.Fatalf("expected an error, got none")
	}

	expErrMsg := "Get \"" + server.URL + "\": context deadline exceeded (Client.Timeout exceeded while awaiting headers)"
	if err.Error() != expErrMsg {
		t.Errorf("expected error message %s, got %s", expErrMsg, err)
	}
}
