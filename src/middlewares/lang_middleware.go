package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"go-clean/src/app/config"
)

func LangMiddleware() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) error {
		lang := fiberCtx.Get("lang", config.GetConfig().App.Locale)
		if fiberCtx.Get("lang") == "" {
			fiberCtx.Request().Header.Add("lang", lang)
			fiberCtx.Locals("lang", lang)
		}
		return fiberCtx.Next()
	}
}
