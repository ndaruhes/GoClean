package main

import (
	"errors"
	users "go-clean/domains/users/handler/http"
	"go-clean/middlewares"
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
			ctx.JSON(http.StatusInternalServerError, responses.BasicResponse{
				Message: "Error Migrate",
				Error:   err,
			})
		} else {
			ctx.JSON(http.StatusOK, responses.BasicResponse{
				Message: "Success Migrate",
			})
		}
	} else {
		ctx.JSON(http.StatusUnauthorized, responses.BasicResponse{
			Message: "Error Migrate (Unauthorized)",
			Error:   errors.New("Unauthorized"),
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
