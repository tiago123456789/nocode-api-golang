package main

import (
	"log"
	"os"

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
)

var actionsBeforePersist map[string]types.ActionInterface

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
	cache := config.GetCache()
	if err != nil {
		log.Fatal(err)
	}

	tableRepository := repository.TableRepositoryNew(db)
	endpointRepository := repository.EndpointRepositoryNew(db)
	customEndpointRepository := repository.CustomEndpointRepositoryNew(db)
	authRespository := repository.AuthRepositoryNew(db)
	authService := service.AuthServiceNew(authRespository)
	tableService := service.TableServiceNew(tableRepository)
	endpointService := service.EndpointServiceNew(tableService, endpointRepository)
	customEndpointService := service.CustomEndpointServiceNew(customEndpointRepository)
	authController := controller.AuthControllerNew(
		*authService,
	)
	tableControler := controller.TableControllerNew(
		*tableService,
	)
	endpointController := controller.EndpointControllerNew(
		*endpointService,
		cache,
	)
	customEndpointController := controller.CustomEndpointControllerNew(
		*customEndpointService,
		actionsBeforePersist,
	)

	err = endpointService.Setup()
	if err != nil {
		log.Fatal(err)
	}

	endpointsFromDB, err := endpointService.GetAllCreated()
	if err != nil {
		log.Fatal(err)
	}

	utils.SetEndpointsInCache(endpointsFromDB)

	app.Use(cors.New())

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
