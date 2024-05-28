package api

import (
	"github.com/andre2ar/zip-temperature/internal/entity"
	"github.com/andre2ar/zip-temperature/internal/services"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"regexp"
	"strings"
)

func GetTemperature(app *entity.App) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		app.Ctx = ctx.UserContext()

		zipcode := ctx.Params("zipcode")
		zipcode = strings.Replace(zipcode, "-", "", -1)

		regex := regexp.MustCompile(`^\d{8}$`)
		if !regex.MatchString(zipcode) {
			_ = ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "invalid zipcode",
			})
			return nil
		}

		_, span := app.Tracer.Start(app.Ctx, "request_get_temperature", oteltrace.WithAttributes(attribute.String("zipcode", zipcode)))
		defer span.End()

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
