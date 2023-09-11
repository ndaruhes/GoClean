package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	server "go-clean/src/app"
	appConfig "go-clean/src/app/config"
)

func main() {
	app := fiber.New()
	server.RegisterMiddlewares(app)
	server.RegisterRoutes(app)
	if err := app.Listen(fmt.Sprintf(":%d", appConfig.GetConfig().App.Port)); err != nil {
		return
	}
}
