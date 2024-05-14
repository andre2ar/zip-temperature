package services

import (
	"github.com/andre2ar/zip-temperature/internal/entity"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTemperatures(t *testing.T) {
	// Mocking a successful request to ViaCep
	mockViaCep := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Localidade": "Example City"}`))
	}))
	defer mockViaCep.Close()

	// Mocking a successful request to Weather API
	mockWeatherAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Current": {"TempC": 20, "TempF": 68}}`))
	}))
	defer mockWeatherAPI.Close()

	app := &entity.App{WeatherApiKey: "test-key"}

	// Testing with mocked servers
	result, err := GetTemperatures(app, "01153001")
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetTemperatures_ErrorGettingViaCep(t *testing.T) {
	app := &entity.App{WeatherApiKey: "test-key"}

	// Testing error case when getting ViaCep data
	result, err := GetTemperatures(app, "invalid-zipcode")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTemperatures_ErrorGettingWeatherAPI(t *testing.T) {
	// Mocking a successful request to ViaCep
	mockViaCep := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Localidade": "Example City"}`))
	}))
	defer mockViaCep.Close()

	app := &entity.App{WeatherApiKey: "test-key"}

	// Testing error case when getting Weather API data
	result, err := GetTemperatures(app, "12345")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTemperatures_ErrorGettingRequest(t *testing.T) {
	// Mocking a failing request
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	app := &entity.App{WeatherApiKey: "test-key"}

	// Testing error case when the request fails
	result, err := GetTemperatures(app, "12345")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTemperatures_BadRequestFromServer(t *testing.T) {
	// Mocking a server returning bad request
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer mockServer.Close()

	app := &entity.App{WeatherApiKey: "test-key"}

	// Testing error case when server returns bad request
	result, err := GetTemperatures(app, "12345")
	assert.Error(t, err)
	assert.Nil(t, result)

	assert.True(t, err.Error() == "bad request")
}
