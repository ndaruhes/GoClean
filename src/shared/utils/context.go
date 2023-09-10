package utils

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

func GetContext(fiberCtx *fiber.Ctx) context.Context {
	originalContext := fiberCtx.Context()
	fiberCtxValues := originalContext.Value("fiberCtx")
	ctx := context.WithValue(originalContext, "fiberCtx", fiberCtxValues)

	return ctx
}
