package http

import (
	"github.com/gofiber/fiber/v2"
	"go-clean/src/app/infrastructures"
	"go-clean/src/domains/auth"
	authRepository "go-clean/src/domains/auth/repositories"
	authUseCase "go-clean/src/domains/auth/usecases"
	"go-clean/src/models/messages"
	"go-clean/src/models/requests"
	"go-clean/src/models/responses"
	"go-clean/src/shared/validators"
	"net/http"
)

type AuthHttp struct {
	authUc auth.AuthUseCase
}

func NewAuthHttp(route *fiber.App) *AuthHttp {
	db := infrastructures.ConnectDatabase()
	authRepo := authRepository.NewAuthRepository(db)
	authUc := authUseCase.NewAuthUseCase(authRepo)

	handler := &AuthHttp{authUc: authUc}
	setRoutes(route, handler)

	return handler
}

func setRoutes(route *fiber.App, handler *AuthHttp) {
	auth := route.Group("auth")
	{
		auth.Post("/register", handler.RegisterWithEmailPassword)
		auth.Post("/login", handler.LoginByPass)
	}
}

func (handler *AuthHttp) RegisterWithEmailPassword(ctx *fiber.Ctx) error {
	request := &requests.RegisterWithEmailPasswordRequest{}
	if err := ctx.BodyParser(&request); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	if formErrors, err := validators.ValidateStruct(ctx, request); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error:      err,
			FormErrors: formErrors,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	err := handler.authUc.RegisterByPass(ctx, request)
	if err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
	} else {
		messages.SendSuccessResponse(ctx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-AUTH-0001",
		})
	}

	return nil
}

func (handler *AuthHttp) LoginByPass(ctx *fiber.Ctx) error {
	request := &requests.LoginRequest{}
	if err := ctx.BodyParser(&request); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	if formErrors, err := validators.ValidateStruct(ctx, request); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error:      err,
			FormErrors: formErrors,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	data, err := handler.authUc.LoginByPass(ctx, request)
	if err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
	} else {
		messages.SendSuccessResponse(ctx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-AUTH-0002",
			Data:        data,
		})
	}

	return nil
}
