package entity

import (
	"github.com/gofiber/fiber/v2"
)

type App struct {
	*fiber.App
	WeatherApiKey string
}
