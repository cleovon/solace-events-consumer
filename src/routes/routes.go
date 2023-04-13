package routes

import (
	"solace-events-consumer/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	app.Get("/healthcheck", controllers.GetHealthCheck)
}
