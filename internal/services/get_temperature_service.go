package services

import (
	"encoding/json"
	"errors"
	"github.com/andre2ar/zip-temperature/internal/dto"
	"github.com/andre2ar/zip-temperature/internal/entity"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"io"
	"log"
	"net/http"
	"net/url"
)

func GetTemperatures(app *entity.App, zipcode string) (*dto.TemperatureResponseDto, error) {
	viaCepResponse, err := getViaCep(app, zipcode)
	if err != nil {
		return nil, err
	}

	weatherApiResponse, err := getWeatherApi(app, viaCepResponse.Localidade, app.WeatherApiKey)
	if err != nil {
		return nil, err
	}

	return &dto.TemperatureResponseDto{
		City:       viaCepResponse.Localidade,
		Celsius:    weatherApiResponse.Current.TempC,
		Fahrenheit: weatherApiResponse.Current.TempF,
		Kevin:      weatherApiResponse.Current.TempC + 273,
	}, nil
}

func getViaCep(app *entity.App, zipcode string) (*dto.ViaCepResponse, error) {
	_, span := app.Tracer.Start(app.Ctx, "get_location_details", oteltrace.WithAttributes(attribute.String("zipcode", zipcode)))
	defer span.End()

	uri := "https://viacep.com.br/ws/" + zipcode + "/json"
	res, err := getRequest(uri)
	if err != nil {
		log.Println(err)
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

func getWeatherApi(app *entity.App, city string, key string) (*dto.WeatherAPIResponse, error) {
	_, span := app.Tracer.Start(app.Ctx, "get_location_weather", oteltrace.WithAttributes(attribute.String("city", city)))
	defer span.End()

	baseUrl, _ := url.Parse("http://api.weatherapi.com/v1/current.json")

	params := url.Values{}
	params.Add("q", city)
	params.Add("key", key)

	baseUrl.RawQuery = params.Encode()

	res, err := getRequest(baseUrl.String())
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

func getRequest(uri string) (*http.Response, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if res.StatusCode == fiber.StatusBadRequest {
		return nil, errors.New("bad request")
	}

	return res, nil
}
