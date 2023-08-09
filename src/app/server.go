package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	users "go-clean/src/domains/auth/handler/http"
	blogs "go-clean/src/domains/blogs/handler/http"
	"go-clean/src/middlewares"
	"go-clean/src/models/messages"
	"go-clean/src/models/responses"
	"go-clean/src/shared/database/migration"
	"go-clean/src/shared/database/seeder"
	"net/http"
	"os"
)

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
	router.GET("/seed-data", seedData)
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

func seedData(ctx *gin.Context) {
	err := seeder.DBSeed()
	if err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		})
	} else {
		messages.SendSuccessResponse(ctx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-DB-0002",
		})
	}
}
