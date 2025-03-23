package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/tiago123456789/nocode-api-golang/internal/config"
	"github.com/tiago123456789/nocode-api-golang/internal/controller"
	"github.com/tiago123456789/nocode-api-golang/internal/middleware"
	"github.com/tiago123456789/nocode-api-golang/internal/repository"
	"github.com/tiago123456789/nocode-api-golang/internal/service"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
	"github.com/tiago123456789/nocode-api-golang/internal/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var actionsBeforePersist map[string]types.ActionInterface

func SetupOtelSDK(
	ctx context.Context,
) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error

		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}

		shutdownFuncs = nil
		return err
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	tracerProvider, err := newTraceProvider(ctx)
	if err != nil {
		handleErr(err)
		return
	}

	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	return
}

func newTraceProvider(ctx context.Context) (*trace.TracerProvider, error) {
	traceExporter, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, err
	}

	// Create resource with service name
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("no-code-api"), // Set application name here
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)

	return traceProvider, nil
}

func main() {
	actionsBeforePersist = map[string]types.ActionInterface{
		"hash": service.HashPasswordActionServiceNew(),
	}

	env := os.Getenv("ENV")
	if env != "production" {
		_ = godotenv.Load()
	}

	app := fiber.New()
	db, err := config.StartDB()

	config.InitCache()
	cache := config.GetCache()
	if err != nil {
		log.Fatal(err)
	}

	logger := config.GetLogger()

	// // traceExporter, err := config.NewTraceExporter()
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }

	// // tracerProvider := config.NewTraceProvider(traceExporter)
	// // defer tracerProvider.Shutdown(context.Background())
	// // otel.SetTracerProvider(tracerProvider)

	// // meterExpoter, err := config.NewMetricExpoter()
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }

	// // meterProvider := config.NewMeterProvider(meterExpoter)
	// // defer meterProvider.Shutdown(context.Background())
	// // otel.SetMeterProvider(meterProvider)

	// // prop := config.NewPropagation()
	// // otel.SetTextMapPropagator(prop)

	// // tracer = otel.Tracer("nocode-api")
	// // meter = otel.Meter("nocode-api")

	// viewCounter, err = meter.Int64Counter(
	// 	"user.views",
	// 	metric.WithDescription("The number of views"),
	// 	metric.WithUnit("{views}"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	ctx := context.Background()
	otelShutdown, err := SetupOtelSDK(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = errors.Join(err, otelShutdown(ctx))
		log.Println(err)
	}()

	tracer := otel.Tracer("no-code-api-golang")

	tableRepository := repository.TableRepositoryNew(db, tracer)
	endpointRepository := repository.EndpointRepositoryNew(db, tracer)
	customEndpointRepository := repository.CustomEndpointRepositoryNew(db)
	authRespository := repository.AuthRepositoryNew(db)
	authService := service.AuthServiceNew(authRespository)
	tableService := service.TableServiceNew(tableRepository, tracer)
	endpointService := service.EndpointServiceNew(tableService, endpointRepository, tracer)
	customEndpointService := service.CustomEndpointServiceNew(customEndpointRepository)
	authController := controller.AuthControllerNew(
		*authService, logger,
	)
	tableControler := controller.TableControllerNew(
		*tableService,
		tracer,
	)
	endpointController := controller.EndpointControllerNew(
		*endpointService,
		cache,
		logger,
		tracer,
	)
	customEndpointController := controller.CustomEndpointControllerNew(
		*customEndpointService,
		actionsBeforePersist,
		logger,
	)

	err = endpointService.Setup()
	if err != nil {
		log.Fatal(err)
	}

	endpointsFromDB, err := endpointService.GetAllCreated(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	utils.SetEndpointsInCache(endpointsFromDB)

	app.Use(cors.New())
	app.Use(otelfiber.Middleware())

	// app.Use(func(c *fiber.Ctx) error {
	// 	ctx, span := tracer.Start(c.Context(), "info")
	// 	defer span.End()
	// 	viewCounter.Add(ctx, 1)

	// 	err := c.Next()
	// 	return err
	// })

	// enabledNewRelic := os.Getenv("NEW_RELIC_ENABLED")
	// if enabledNewRelic == "yes" {
	// 	app.Use(fibernewrelic.New(fibernewrelic.Config{
	// 		License: os.Getenv("NEW_RELIC_LICENSE_KEY"),
	// 		AppName: os.Getenv("NEW_RELIC_APP_NAME"),
	// 		Enabled: true,
	// 	}))
	// }

	app.Post("auth/login", authController.Login)

	app.Get("/tables",
		middleware.HttpLogs,
		middleware.IsInternalAuthorized,
		tableControler.GetAll)

	app.Get(
		"/tables/:table/columns",
		middleware.HttpLogs,
		middleware.IsInternalAuthorized,
		tableControler.GetColumnsFromTable)

	app.Get("/endpoints",
		middleware.HttpLogs,
		middleware.IsInternalAuthorized,
		endpointController.GetAllCreated)

	app.Delete("/endpoints/:id",
		middleware.HttpLogs,
		middleware.IsInternalAuthorized,
		endpointController.DeleteById)

	app.Post("/endpoints",
		middleware.HttpLogs,
		middleware.IsInternalAuthorized,
		endpointController.Create)

	app.Put("/:table/:id",
		middleware.HttpLogs,
		middleware.IsAuthorized(),
		customEndpointController.Put)

	app.Post("/*",
		middleware.HttpLogs,
		middleware.IsAuthorized(),
		customEndpointController.Post)

	app.Delete("/:table/:id",
		middleware.HttpLogs,
		middleware.IsAuthorized(),
		customEndpointController.Delete)

	app.Get("/:table/:id",
		middleware.HttpLogs,
		middleware.IsAuthorized(),
		customEndpointController.GetById)

	app.Get("/*",
		middleware.HttpLogs,
		middleware.IsAuthorized(),
		middleware.CacheResponse(cache),
		customEndpointController.GetAll)

	app.Listen(":3000")
}
