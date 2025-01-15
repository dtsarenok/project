package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/dtsarenok/project/pkg/parser"
)

// MockCalculate заменяет функцию Calculate для тестирования
func MockCalculate(expr string) (float64, error) {
	if expr == "2+2" {
		return 4, nil
	}
	return 0, fmt.Errorf("invalid expression")
}

func TestCalculateHandler(t *testing.T) {
	// Заменяем оригинальную функцию Calculate на MockCalculate
	originalCalculate := parser.Calculate
	defer func() { parser.Calculate = originalCalculate }()
	parser.Calculate = MockCalculate

	tests := []struct {
		name           string
		expression     string
		expectedResult float64
		expectedError  string
		expectedCode   int
	}{
		{
			name:           "Valid expression",
			expression:     "2+2",
			expectedResult: 4,
			expectedError:  "",
			expectedCode:   http.StatusOK,
		},
		{
			name:           "Invalid expression",
			expression:     "invalid",
			expectedResult: 0,
			expectedError:  "Expression is not valid",
			expectedCode:   http.StatusUnprocessableEntity,
		},
		{
			name:           "Invalid JSON format",
			expression:     "",
			expectedResult: 0,
			expectedError:  "Invalid JSON format",
			expectedCode:   http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if tt.name == "Invalid JSON format" {
				body = []byte("invalid json")
			} else {
				req := Request{Expression: tt.expression}
				body, _ = json.Marshal(req)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			calculateHandler(w, req)

			res := w.Result()
			if res.StatusCode != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, res.StatusCode)
			}

			var response Response
			if err := json.NewDecoder(res.Body).Decode(&response); err == nil {
				if response.Error != tt.expectedError {
					t.Errorf("expected error %q, got %q", tt.expectedError, response.Error)
				}
				if response.Error == "" { // Only check result if there is no error
					resultAsFloat, err := strconv.ParseFloat(fmt.Sprintf("%f", response.Result), 64)
					if err != nil {
						t.Errorf("failed to convert result to float64: %v", err)
					}
					if resultAsFloat != tt.expectedResult {
						t.Errorf("expected result %f, got %f", tt.expectedResult, resultAsFloat)
					}
				}
			}
		})
	}
}
