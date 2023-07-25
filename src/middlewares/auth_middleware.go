package middlewares

import (
	"go-clean/src/models/messages"
	"go-clean/src/models/responses"
	"go-clean/src/shared/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := helpers.VerifyToken(ctx)
		if err != nil {
			messages.SendErrorResponse(ctx, responses.ErrorResponse{
				Error:      err,
				StatusCode: http.StatusUnauthorized,
			})
			ctx.Abort()
			return
		}

		if user.Role == "Member" {
			ctx.Set("member", user)
		} else {
			ctx.Set("admin", user)
		}

		ctx.Next()
	}
}
