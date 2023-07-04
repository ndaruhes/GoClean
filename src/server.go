package main

import (
	blogs "go-clean/domains/blogs/handler/http"
	users "go-clean/domains/users/handler/http"
	"go-clean/middlewares"
	"go-clean/models/messages"
	"go-clean/models/responses"
	"net/http"
	"os"

	"go-clean/configs/migration"

	"github.com/gin-gonic/gin"
)

func migrate(ctx *gin.Context) {
	if ctx.Query("key") == os.Getenv("MIGRATE_KEY") {
		err := migration.Migrate()
		if err != nil {
			messages.SendErrorResponse(ctx, responses.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
			})
		} else {
			messages.SendSuccessResponse(ctx, responses.SuccessResponse{
				SuccessCode: "SUCCESS-DB-0001",
			})
		}
	} else {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
		})
	}
}

func RegisterMiddlewares(router *gin.Engine) {
	router.Use(middlewares.LangMiddleware())
}

func RegisterRoutes(router *gin.Engine) {
	router.GET("/", func(ctx *gin.Context) {
		messages.SendSuccessResponse(ctx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-BASIC-0001",
			StatusCode:  http.StatusOK,
		})
	})
	router.Static("/images", "./public/images")
	router.GET("/migrate", migrate)
	users.NewAuthHttp(router)
	blogs.NewBlogHttp(router)
}
