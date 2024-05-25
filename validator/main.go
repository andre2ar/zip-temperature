package main

import (
	"encoding/json"
	"github.com/andre2ar/zip-temperature/validator/dto"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/temperature", temperature)

	log.Println("Listening on localhost:8081")
	log.Fatalln(http.ListenAndServe(":8081", nil))
}

func temperature(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"message":"method not allowed"}`))
		return
	}

	var input dto.Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"Invalid request"}`))
		return
	}

	if len(input.Zipcode) != 8 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"message":"Invalid zipcode"}`))
		return
	}

	zipTemperatureResponse, err := getTemperature(input)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"Failed to get temperature"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(zipTemperatureResponse)
	w.Write(jsonData)
}

func getTemperature(input dto.Input) (*dto.ZipTemperatureResponse, error) {
	res, err := http.Get("http://host.docker.internal:8080/api/v1/temperature/" + input.Zipcode)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	responseBody, _ := io.ReadAll(res.Body)
	var zipTemperatureResponse dto.ZipTemperatureResponse
	err = json.Unmarshal(responseBody, &zipTemperatureResponse)
	if err != nil {
		return nil, err
	}

	return &zipTemperatureResponse, nil
}
