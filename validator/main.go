package main

import (
	"context"
	"encoding/json"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"io"
	"log"
	"net/http"
)

var tracer = otel.Tracer("zipcode_validator")

func main() {
	tp := InitTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	app := fiber.New()

	app.Use(logger.New())
	app.Use(requestid.New())
	app.Use(otelfiber.Middleware())

	app.Post("/api/v1/temperature", temperature)

	log.Fatal(app.Listen(":8081"))
}

func InitTracer() *sdktrace.TracerProvider {
	// Will take OTEL_EXPORTER_OTLP_ENDPOINT environment variable as the collector path
	clientOTel := otlptracegrpc.NewClient()
	exporter, err := otlptrace.New(context.Background(), clientOTel)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %e", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("validator"),
			)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

type ZipTemperatureResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type Input struct {
	Zipcode string `json:"zipcode"`
}

func temperature(c *fiber.Ctx) error {
	payload := Input{}
	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request",
		})
	}

	_, span := tracer.Start(c.UserContext(), "request_temperature", oteltrace.WithAttributes(attribute.String("zipcode", payload.Zipcode)))
	defer span.End()

	if !isValidZipcode(c.UserContext(), payload.Zipcode) {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "invalid zipcode",
		})
	}

	zipTemperatureResponse, err := getTemperature(c.UserContext(), payload)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get temperature",
		})
	}

	return c.JSON(zipTemperatureResponse)
}

func getTemperature(ctx context.Context, input Input) (*ZipTemperatureResponse, error) {
	_, span := tracer.Start(ctx, "get_temperature", oteltrace.WithAttributes(attribute.String("zipcode", input.Zipcode)))
	defer span.End()

	res, err := http.Get("http://host.docker.internal:8080/api/v1/temperature/" + input.Zipcode)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	responseBody, _ := io.ReadAll(res.Body)
	var zipTemperatureResponse ZipTemperatureResponse
	err = json.Unmarshal(responseBody, &zipTemperatureResponse)
	if err != nil {
		return nil, err
	}

	return &zipTemperatureResponse, nil
}

func isValidZipcode(ctx context.Context, zipcode string) bool {
	_, span := tracer.Start(ctx, "validate_zipcode", oteltrace.WithAttributes(attribute.String("zipcode", zipcode)))
	defer span.End()

	if len(zipcode) != 8 {
		return false
	}

	return true
}
