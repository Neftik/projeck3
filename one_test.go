package main

import (
	"bytes"
	//"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer([]byte(`{"expression": "2+2*2"}`)))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	http.HandlerFunc(calculateHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK && rr.Code != http.StatusCreated {
		t.Errorf("expected status 201 or 200, got %v", rr.Code)
	}
} 

func TestGetExpressionsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/expressions", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	http.HandlerFunc(getExpressionsHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %v", rr.Code)
	}
}

func TestGetTaskHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/internal/task", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	http.HandlerFunc(getTaskHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK && rr.Code != http.StatusNotFound {
		t.Errorf("expected status 200 or 404, got %v", rr.Code)
	}
}

func TestSubmitResultHandler(t *testing.T) {
	payload := bytes.NewBuffer([]byte(`{"id": "task_1", "result": 4.0}`))
	req, err := http.NewRequest("POST", "/internal/task", payload)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	http.HandlerFunc(submitResultHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %v", rr.Code)
	}
}
