package main

import (
	"context"
	"errors"
	"github.com/andre2ar/zip-temperature/config"
	"github.com/andre2ar/zip-temperature/internal/infra/rest"
	"github.com/andre2ar/zip-temperature/pkg"
	"go.opentelemetry.io/otel"
	"log"
)

var tracer = otel.Tracer("github.com/andre2ar/zip-temperature")

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	if cfg.WeatherAPIKey == "" {
		err = errors.New("Weather API key is required")
		panic(err)
	}

	tp := pkg.InitTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	log.Fatalln(rest.CreateRestServer(cfg, tracer))
}
