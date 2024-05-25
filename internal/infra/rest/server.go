package rest

import (
	"github.com/andre2ar/zip-temperature/config"
	"github.com/andre2ar/zip-temperature/internal/entity"
	"github.com/andre2ar/zip-temperature/internal/infra/routes"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.opentelemetry.io/otel/trace"
)

func CreateRestServer(cfg *config.Config, tracer trace.Tracer) error {
	app := entity.App{
		App:           fiber.New(),
		WeatherApiKey: cfg.WeatherAPIKey,
		Tracer:        tracer,
	}

	app.Use(logger.New())
	app.Use(requestid.New())
	app.Use(otelfiber.Middleware())

	api := app.Group("/api")
	apiV1 := api.Group("/v1")
	routes.RegisterAPI(apiV1, &app)

	err := app.Listen(cfg.WebServerPort)
	if err != nil {
		return err
	}

	return nil
}
