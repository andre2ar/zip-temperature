package main

import (
	"errors"
	"github.com/andre2ar/zip-temperature/config"
	"github.com/andre2ar/zip-temperature/internal/infra/rest"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	if cfg.WeatherAPIKey == "" {
		err = errors.New("Weather API key is required")
		panic(err)
	}

	log.Fatalln(rest.CreateRestServer(cfg))
}
