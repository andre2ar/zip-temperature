package api

import (
	"github.com/andre2ar/zip-temperature/internal/entity"
	"github.com/andre2ar/zip-temperature/internal/services"
	"github.com/gofiber/fiber/v2"
	"regexp"
	"strings"
)

func GetTemperature(app *entity.App) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		zipcode := ctx.Params("zipcode")
		zipcode = strings.Replace(zipcode, "-", "", -1)

		regex := regexp.MustCompile(`^\d{8}$`)
		if !regex.MatchString(zipcode) {
			_ = ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "invalid zipcode",
			})
			return nil
		}

		temperatures, err := services.GetTemperatures(app, zipcode)
		if err != nil {
			_ = ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "can not find zipcode",
			})
			return nil
		}

		_ = ctx.Status(fiber.StatusOK).JSON(temperatures)

		return nil
	}
}
