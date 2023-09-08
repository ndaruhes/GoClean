package utils

import (
	"github.com/gofiber/fiber/v2"
)

func OverrideFiberRequest(ctx *fiber.Ctx, key string, value interface{}) {
	ctx.Locals(key, value)
}
