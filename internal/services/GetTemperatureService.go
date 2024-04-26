package services

import (
	"encoding/json"
	"errors"
	"github.com/andre2ar/zip-temperature/internal/dto"
	"github.com/andre2ar/zip-temperature/internal/entity"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"net/http"
)

func GetTemperatures(app *entity.App, zipcode string) (*dto.TemperatureResponseDto, error) {
	viaCepResponse, err := getViaCep(zipcode)
	if err != nil {
		return nil, err
	}

	weatherApiResponse, err := getWeatherApi(viaCepResponse.Localidade, app.WeatherApiKey)
	if err != nil {
		return nil, err
	}

	return &dto.TemperatureResponseDto{
		Celsius:    weatherApiResponse.Current.TempC,
		Fahrenheit: weatherApiResponse.Current.TempF,
		Kevin:      weatherApiResponse.Current.TempC + 273,
	}, nil
}

func getViaCep(zipcode string) (*dto.ViaCepResponse, error) {
	url := "https://viacep.com.br/ws/" + zipcode + "/json"
	res, err := getRequest(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	responseBody, _ := io.ReadAll(res.Body)
	var viaCepResponse dto.ViaCepResponse
	err = json.Unmarshal(responseBody, &viaCepResponse)
	if err != nil {
		return nil, err
	}

	return &viaCepResponse, nil
}

func getWeatherApi(city string, key string) (*dto.WeatherAPIResponse, error) {
	url := "http://api.weatherapi.com/v1/current.json?key=" + key + "&q=" + city
	res, err := getRequest(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	responseBody, _ := io.ReadAll(res.Body)
	var weatherAPIResponse dto.WeatherAPIResponse
	err = json.Unmarshal(responseBody, &weatherAPIResponse)
	if err != nil {
		return nil, err
	}

	return &weatherAPIResponse, nil
}

func getRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if res.StatusCode == fiber.StatusBadRequest {
		return nil, errors.New("Bad Request")
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return res, nil
}
