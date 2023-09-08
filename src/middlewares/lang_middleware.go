package middlewares

import (
	"go-clean/src/app/config"
	"go-clean/src/shared/utils"

	"github.com/gofiber/fiber/v2"
)

func LangMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Get("lang") == "" {
			lang := c.Get("lang", config.GetConfig().App.Locale)
			c.Locals("lang", lang)
			utils.OverrideFiberRequest(c, "lang", lang)
		}
		return c.Next()
	}
}
