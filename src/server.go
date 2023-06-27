package main

import (
	users "go-clean/domains/users/handler/http"
	"go-clean/middlewares"
	"go-clean/models/messages"
	"go-clean/models/responses"
	"os"

	"go-clean/configs/migration"

	"github.com/gin-gonic/gin"
)

func migrate(ctx *gin.Context) {
	if ctx.Query("key") == os.Getenv("MIGRATE_KEY") {
		err := migration.Migrate()
		if err != nil {
			messages.SendErrorResponse(ctx, responses.ErrorResponse{
				ErrorCode: "ERROR-50001",
			})
		} else {
			messages.SendSuccessResponse(ctx, responses.SuccessResponse{
				SuccessCode: "SUCCESS-0003",
			})
		}
	} else {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			ErrorCode: "ERROR-50002",
		})
	}
}

func RegisterMiddlewares(router *gin.Engine) {
	router.Use(middlewares.LangMiddleware())
}

func RegisterRoutes(router *gin.Engine) {
	router.GET("/migrate", migrate)
	users.NewAuthHttp(router)
}
