package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsValidCEPFormat(t *testing.T) {
	tests := []struct {
		name     string
		cep      string
		expected bool
	}{
		{"Valid CEP - 8 digits", "01310100", true},
		{"Valid CEP - with hyphen", "01310-100", true},
		{"Invalid CEP - 7 digits", "0131010", false},
		{"Invalid CEP - 9 digits", "013101000", false},
		{"Invalid CEP - with letters", "0131010a", false},
		{"Invalid CEP - empty", "", false},
		{"Invalid CEP - special chars", "01310-10@", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidCEPFormat(tt.cep)
			if result != tt.expected {
				t.Errorf("isValidCEPFormat(%s) = %v; expected %v", tt.cep, result, tt.expected)
			}
		})
	}
}

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{"Zero Celsius", 0, 32},
		{"Positive temperature", 28.5, 83.3},
		{"Negative temperature", -10, 14},
		{"Boiling point", 100, 212},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := celsiusToFahrenheit(tt.celsius)
			if result != tt.expected {
				t.Errorf("celsiusToFahrenheit(%f) = %f; expected %f", tt.celsius, result, tt.expected)
			}
		})
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{"Zero Celsius", 0, 273},
		{"Positive temperature", 28.5, 301.5},
		{"Negative temperature", -10, 263},
		{"Absolute zero", -273, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := celsiusToKelvin(tt.celsius)
			if result != tt.expected {
				t.Errorf("celsiusToKelvin(%f) = %f; expected %f", tt.celsius, result, tt.expected)
			}
		})
	}
}

func TestHandleWeatherInvalidFormat(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/weather/123", nil)
	w := httptest.NewRecorder()

	handleWeather(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}

	var response ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Message != "invalid zipcode" {
		t.Errorf("Expected message 'invalid zipcode', got '%s'", response.Message)
	}
}

func TestHandleWeatherInvalidCEP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/weather/99999999", nil)
	w := httptest.NewRecorder()

	handleWeather(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	var response ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Message != "can not find zipcode" {
		t.Errorf("Expected message 'can not find zipcode', got '%s'", response.Message)
	}
}

func TestHandleHealthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", w.Body.String())
	}
}

func TestHandleWeatherMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/weather/01310100", nil)
	w := httptest.NewRecorder()

	handleWeather(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestGetLocationByCEP(t *testing.T) {
	// Teste com CEP válido da Avenida Paulista
	location, err := getLocationByCEP("01310100")

	if err != nil {
		t.Errorf("Expected no error for valid CEP, got: %v", err)
	}

	if location == "" {
		t.Error("Expected location to be returned, got empty string")
	}

	// Teste com CEP inválido
	_, err = getLocationByCEP("99999999")

	if err == nil {
		t.Error("Expected error for invalid CEP, got nil")
	}
}

func TestRespondWithJSON(t *testing.T) {
	w := httptest.NewRecorder()

	payload := TemperatureResponse{
		TempC: 25.0,
		TempF: 77.0,
		TempK: 298.0,
	}

	respondWithJSON(w, http.StatusOK, payload)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	var response TemperatureResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.TempC != payload.TempC {
		t.Errorf("Expected TempC %f, got %f", payload.TempC, response.TempC)
	}
}
