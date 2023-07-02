package middlewares

import (
	"github.com/gin-gonic/gin"
	"go-clean/models/messages"
	"go-clean/models/responses"
	"go-clean/shared/helpers"
	"net/http"
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
