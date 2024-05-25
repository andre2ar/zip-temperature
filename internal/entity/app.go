package entity

import (
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

type App struct {
	*fiber.App
	WeatherApiKey string
	Tracer        trace.Tracer
}
