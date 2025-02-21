package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/tiago123456789/nocode-api-golang/internal/config"
	"github.com/tiago123456789/nocode-api-golang/internal/middleware"
	"github.com/tiago123456789/nocode-api-golang/internal/repository"
	"github.com/tiago123456789/nocode-api-golang/internal/service"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
	"github.com/tiago123456789/nocode-api-golang/internal/utils"
)

var endpoints map[string]types.Endpoint
var actionsBeforePersist map[string]types.ActionInterface

func main() {
	endpoints = map[string]types.Endpoint{}
	actionsBeforePersist = map[string]types.ActionInterface{
		"hash": service.HashPasswordActionServiceNew(),
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	db, err := config.StartDB()

	if err != nil {
		log.Fatal(err)
	}

	tableRepository := repository.TableRepositoryNew(db)
	endpointRepository := repository.EndpointRepositoryNew(db)
	customEndpointRepository := repository.CustomEndpointRepositoryNew(db)
	tableService := service.TableServiceNew(tableRepository)
	endpointService := service.EndpointServiceNew(db, tableService, endpointRepository)
	customEndpoint := service.CustomEndpointServiceNew(db, customEndpointRepository)

	err = endpointService.Setup()
	if err != nil {
		log.Fatal(err)
	}

	endpointsFromDB, err := endpointService.GetAllCreated()
	if err != nil {
		log.Fatal(err)
	}

	endpoints = endpointsFromDB

	app.Get("/tables",
		middleware.HttpLogs,
		middleware.IsInternalAuthorized,
		func(c *fiber.Ctx) error {
			results, _ := tableService.GetAll()
			return c.JSON(fiber.Map{
				"data": results,
			})
		})

	app.Get(
		"/tables/:table/columns",
		middleware.HttpLogs,
		middleware.IsInternalAuthorized,
		func(c *fiber.Ctx) error {
			results, _ := tableService.GetColumnsFromTable(c.Params("table"))
			return c.JSON(fiber.Map{
				"data": results,
			})
		})

	app.Get("/endpoints",
		middleware.HttpLogs,
		middleware.IsInternalAuthorized,
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"data": endpoints,
			})
		})

	app.Post("/endpoints",
		middleware.HttpLogs,
		middleware.IsInternalAuthorized,
		func(c *fiber.Ctx) error {
			var endpoint types.Endpoint
			c.BodyParser(&endpoint)

			_, err := endpointService.Create(endpoint)
			if err != nil {
				return c.Status(409).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			endpoints[endpoint.Path] = endpoint
			return c.SendStatus(201)
		})

	app.Put("/:table/:id",
		middleware.HttpLogs,
		middleware.IsAuthorized(endpoints),
		func(c *fiber.Ctx) error {
			newRegister := map[string]interface{}{}
			c.BodyParser(&newRegister)

			err := customEndpoint.Put(
				newRegister, c.Params("table"), c.Params("id"),
			)
			if err != nil {
				return c.Status(404).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			return c.JSON(fiber.Map{
				"message": "ok",
				"id":      c.Params("id"),
			})
		})

	app.Post("/*",
		middleware.HttpLogs,
		middleware.IsAuthorized(endpoints),
		func(c *fiber.Ctx) error {
			endpoint := endpoints[c.Path()]
			newRegister := map[string]interface{}{}
			c.BodyParser(&newRegister)

			if endpoint.Validations != nil || len(endpoint.Validations) > 0 {
				validationErrors := []string{}
				for _, value := range endpoint.Validations {
					err := utils.IsValid(
						newRegister[value.Field],
						value.Rules,
						value.Field,
					)

					if err != nil {
						validationErrors = append(validationErrors, err.Error())
					}
				}

				if len(validationErrors) > 0 {
					return c.Status(400).JSON(fiber.Map{
						"error": validationErrors,
					})
				}
			}

			if len(endpoint.ActionsBeforePersist) > 0 {
				for _, item := range endpoint.ActionsBeforePersist {
					newRegister[item.Field] = actionsBeforePersist[item.Action].Apply(
						newRegister[item.Field],
					)
				}
			}

			id, err := customEndpoint.Post(newRegister, endpoint.Table)
			if err != nil {
				fmt.Print(err)
				return c.Status(500).JSON(fiber.Map{
					"message": "Interval server error",
				})
			}

			return c.JSON(fiber.Map{
				"message": "ok",
				"id":      id,
			})
		})

	app.Delete("/:table/:id",
		middleware.HttpLogs,
		middleware.IsAuthorized(endpoints),
		func(c *fiber.Ctx) error {
			err := customEndpoint.Delete(c.Params("table"), c.Params("id"))

			if err != nil {
				return c.Status(404).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			return c.SendStatus(204)
		})

	app.Get("/:table/:id",
		middleware.HttpLogs,
		middleware.IsAuthorized(endpoints),
		func(c *fiber.Ctx) error {
			results, _ := customEndpoint.GetById(c.Params("table"), c.Params("id"))
			if len(results) == 0 {
				return c.SendStatus(404)
			}

			return c.JSON(fiber.Map{
				"data": results[0],
			})
		})

	app.Get("/*",
		middleware.HttpLogs,
		middleware.IsAuthorized(endpoints),
		func(c *fiber.Ctx) error {
			endpoint := endpoints[c.Path()]
			var params []interface{}

			if endpoint.Query != "" {
				for _, value := range endpoint.QueryParams {
					queryStringValue := c.Query(value)
					if queryStringValue != "" {
						params = append(params, queryStringValue)
					} else if c.Params(value) != "" {
						params = append(params, c.Params(value))
					}
				}

				if len(params) != len(endpoint.QueryParams) {
					return c.Status(400).JSON(fiber.Map{
						"error": fmt.Sprintf(
							"You need to provide the following params via querystring: %s",
							strings.Join(endpoint.QueryParams, ","),
						),
					})
				}

				results, _ := customEndpoint.Get(endpoint, params)
				return c.JSON(fiber.Map{
					"data": results,
				})
			}

			results, _ := customEndpoint.Get(endpoint, params)
			return c.JSON(fiber.Map{
				"data": results,
			})
		})

	app.Listen(":3000")
}
