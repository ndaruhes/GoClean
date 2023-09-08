package middlewares

import (
	"go-clean/src/models/messages"
	"go-clean/src/models/responses"
	"go-clean/src/shared/helpers"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func Authenticated() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user, err := helpers.VerifyToken(ctx)
		if err != nil {
			messages.SendErrorResponse(ctx, responses.ErrorResponse{
				Error:      err,
				StatusCode: http.StatusUnauthorized,
			})
			return nil
		}

		if user.Role == "Member" {
			ctx.Locals("member", user)
		} else {
			ctx.Locals("admin", user)
		}

		return ctx.Next()
	}
}
