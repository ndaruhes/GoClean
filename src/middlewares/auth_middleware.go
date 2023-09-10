package middlewares

import (
	"go-clean/src/models/messages"
	"go-clean/src/models/responses"
	"go-clean/src/shared/helpers"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func Authenticated() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) error {
		user, err := helpers.VerifyToken(fiberCtx)
		if err != nil {
			messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
				Error:      err,
				StatusCode: http.StatusUnauthorized,
			})
			return nil
		}

		if user.Role == "Member" {
			fiberCtx.Locals("member", user)
		} else {
			fiberCtx.Locals("admin", user)
		}

		return fiberCtx.Next()
	}
}
