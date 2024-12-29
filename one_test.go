package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCalcFunction(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
		hasError   bool
	}{
		{"2+2", 4, false},
		{"2-3", -1, false},
		{"2*3", 6, false},
		{"6/3", 2, false},
		{"(2+3)*4", 20, false},
		{"2+", 0, true},
		{"2/0", 0, true},
		{"invalid", 0, true},
	}

	for _, test := range tests {
		result, err := Calc(test.expression)
		if test.hasError {
			if err == nil {
				t.Errorf("expected an error for expression %q, but got none", test.expression)
			}
		} else {
			if err != nil {
				t.Errorf("did not expect an error for expression %q, but got %v", test.expression, err)
			}
			if result != test.expected {
				t.Errorf("for expression %q, expected %f but got %f", test.expression, test.expected, result)
			}
		}
	}
}

func TestTokenizeFunction(t *testing.T) {
	expression := "(2+3)*4"
	expected := []string{"(", "2", "+", "3", ")", "*", "4"}

	tokens, err := tokenize(expression)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, token := range tokens {
		if token != expected[i] {
			t.Errorf("expected token %q but got %q", expected[i], token)
		}
	}
}

func TestCalculateHandlerFunction(t *testing.T) {
	tests := []struct {
		requestBody  string
		expectedCode int
		expectedBody map[string]string
	}{
		{"{\"expression\": \"2+2\"}", http.StatusOK, map[string]string{"result": "4"}},
		{"{\"expression\": \"2/0\"}", http.StatusUnprocessableEntity, map[string]string{"error": "evaluation error: division by zero"}},
		{"{\"expression\": \"\"}", http.StatusBadRequest, map[string]string{"error": "Expression cannot be empty"}},
		{"", http.StatusBadRequest, map[string]string{"error": "Invalid request body"}},
	}

	for _, test := range tests {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer([]byte(test.requestBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		calculateHandler(w, req)

		resp := w.Result()
		if resp.StatusCode != test.expectedCode {
			t.Errorf("expected status code %d but got %d", test.expectedCode, resp.StatusCode)
		}

		var resBody map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&resBody); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

		if !reflect.DeepEqual(resBody, test.expectedBody) {
			t.Errorf("expected body %v but got %v", test.expectedBody, resBody)
		}
	}
}
