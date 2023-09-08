package utils

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

func OverrideFiberRequest(ctx *fiber.Ctx, key string, value interface{}) {
	c := context.WithValue(ctx.Context(), key, value)
	ctx.Locals(c)
}
