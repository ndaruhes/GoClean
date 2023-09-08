package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"go-clean/src/app/config"
	"go-clean/src/shared/utils"
)

func LangMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if ctx.Get("lang") == "" {
			lang := ctx.Get("lang", config.GetConfig().App.Locale)
			ctx.Locals("lang", lang)
			utils.OverrideFiberRequest(ctx, "lang", lang)
		}
		return ctx.Next()
	}
}
