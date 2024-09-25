package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"go-clean/src/app/config"
)

func LangMiddleware() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) error {
		langHeader := fiberCtx.Get("lang")
		if langHeader == "" {
			langHeader = config.GetConfig().App.DefaultLocale
		}

		fiberCtx.Request().Header.Add("lang", langHeader)
		fiberCtx.Locals("lang", langHeader)

		return fiberCtx.Next()
	}
}
