package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	users "go-clean/src/domains/auth/handler/http"
	blogs "go-clean/src/domains/blogs/handler/http"
	"go-clean/src/middlewares"
	"go-clean/src/models/messages"
	"go-clean/src/models/responses"
	appConfig "go-clean/src/setup/app"
	"go-clean/src/setup/migration"
	"net/http"
	"os"
)

func main() {
	app := gin.Default()
	RegisterMiddlewares(app)
	RegisterRoutes(app)
	err := app.Run(fmt.Sprintf(":%d", appConfig.GetConfig().App.Port))
	if err != nil {
		return
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

func migrate(ctx *gin.Context) {
	if ctx.Query("key") == os.Getenv("MIGRATE_KEY") {
		err := migration.Migrate()
		if err != nil {
			messages.SendErrorResponse(ctx, responses.ErrorResponse{
				Error:      err,
				StatusCode: http.StatusInternalServerError,
			})
		} else {
			messages.SendSuccessResponse(ctx, responses.SuccessResponse{
				SuccessCode: "SUCCESS-DB-0001",
			})
		}
	} else {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error:      errors.New("key is invalid"),
			StatusCode: http.StatusInternalServerError,
		})
	}
}
