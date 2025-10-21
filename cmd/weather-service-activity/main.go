package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
	Erro       bool   `json:"erro,omitempty"`
}

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type TemperatureResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func main() {
	// Carrega variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/weather/", handleWeather)
	http.HandleFunc("/health", handleHealth)

	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func handleWeather(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extrai o CEP da URL
	path := strings.TrimPrefix(r.URL.Path, "/weather/")
	cep := strings.TrimSpace(path)

	// Valida o formato do CEP
	if !isValidCEPFormat(cep) {
		respondWithError(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	// Busca a localização pelo CEP
	location, err := getLocationByCEP(cep)
	if err != nil {
		if err.Error() == "zipcode not found" {
			respondWithError(w, http.StatusNotFound, "can not find zipcode")
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	// Busca a temperatura da localização
	tempC, err := getTemperature(location)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error fetching weather data")
		return
	}

	// Calcula as conversões
	response := TemperatureResponse{
		TempC: tempC,
		TempF: celsiusToFahrenheit(tempC),
		TempK: celsiusToKelvin(tempC),
	}

	respondWithJSON(w, http.StatusOK, response)
}

func isValidCEPFormat(cep string) bool {
	// Remove hífens se existirem
	cep = strings.ReplaceAll(cep, "-", "")

	// Valida se tem exatamente 8 dígitos numéricos
	matched, _ := regexp.MatchString(`^\d{8}$`, cep)
	return matched
}

func getLocationByCEP(cep string) (string, error) {
	// Remove hífens para padronizar
	cep = strings.ReplaceAll(cep, "-", "")

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var viaCEP ViaCEPResponse
	if err := json.Unmarshal(body, &viaCEP); err != nil {
		return "", err
	}

	// Verifica se o CEP foi encontrado
	if viaCEP.Erro || viaCEP.Localidade == "" {
		return "", fmt.Errorf("zipcode not found")
	}

	return viaCEP.Localidade, nil
}

func getTemperature(location string) (float64, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		apiKey = "demo" // Para testes locais
	}

	encodedLocation := url.QueryEscape(location)
	apiURL := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, encodedLocation)

	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var weatherData WeatherAPIResponse
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return 0, err
	}

	return weatherData.Current.TempC, nil
}

func celsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

func celsiusToKelvin(celsius float64) float64 {
	return celsius + 273
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{Message: message})
}
