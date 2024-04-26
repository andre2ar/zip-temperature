package routes

import (
	"github.com/andre2ar/zip-temperature/internal/entity"
	"github.com/andre2ar/zip-temperature/internal/infra/controllers/api"
	"github.com/gofiber/fiber/v2"
)

func RegisterAPI(router fiber.Router, app *entity.App) {
	router.Get("/temperature/:zipcode", api.GetTemperature(app))
}
