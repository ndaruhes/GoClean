package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go-clean/src/app/config"
)

func CorsMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     config.GetConfig().Cors.AllowOrigins,
		AllowHeaders:     config.GetConfig().Cors.AllowHeaders,
		AllowMethods:     config.GetConfig().Cors.AllowMethods,
		AllowCredentials: config.GetConfig().Cors.AllowCredentials,
		ExposeHeaders:    config.GetConfig().Cors.ExposeHeaders,
	})
}
