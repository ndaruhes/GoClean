package server

import (
	"errors"
	users "go-clean/src/domains/auth/handler/http"
	blogs "go-clean/src/domains/blogs/handler/http"
	"go-clean/src/middlewares"
	"go-clean/src/models/messages"
	"go-clean/src/models/responses"
	"go-clean/src/shared/database/migration"
	"go-clean/src/shared/database/seeder"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

func RegisterMiddlewares(router *fiber.App) {
	router.Use(middlewares.LangMiddleware())
}

func RegisterRoutes(router *fiber.App) {
	router.Get("/", func(ctx *fiber.Ctx) error {
		messages.SendSuccessResponse(ctx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-BASIC-0001",
			StatusCode:  http.StatusOK,
		})

		return nil
	})
	router.Static("/images", "./public/images")
	router.Get("/migrate", migrate)
	router.Get("/seed-data", seedData)
	users.NewAuthHttp(router)
	blogs.NewBlogHttp(router)
}

func migrate(ctx *fiber.Ctx) error {
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

	return nil
}

func seedData(ctx *fiber.Ctx) error {
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

	return nil
}
