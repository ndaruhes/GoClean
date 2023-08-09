package http

import (
	"go-clean/src/app/infrastructures"
	"go-clean/src/domains/auth"
	authRepository "go-clean/src/domains/auth/repositories"
	authUseCase "go-clean/src/domains/auth/usecases"
	"go-clean/src/models/messages"
	"go-clean/src/models/requests"
	"go-clean/src/models/responses"
	"go-clean/src/shared/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHttp struct {
	authUc auth.AuthUseCase
}

func NewAuthHttp(route *gin.Engine) *AuthHttp {
	db := infrastructures.ConnectDatabase()
	authRepo := authRepository.NewAuthRepository(db)
	authUc := authUseCase.NewAuthUseCase(authRepo)

	handler := &AuthHttp{authUc: authUc}
	setRoutes(route, handler)

	return handler
}

func setRoutes(route *gin.Engine, handler *AuthHttp) {
	auth := route.Group("auth")
	{
		auth.POST("/register", handler.RegisterWithEmailPassword)
		auth.POST("/login", handler.LoginByPass)
	}
}

func (handler *AuthHttp) RegisterWithEmailPassword(ctx *gin.Context) {
	request := &requests.RegisterWithEmailPasswordRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	if formErrors, err := validators.ValidateStruct(ctx, request); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error:      err,
			FormErrors: formErrors,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	err := handler.authUc.RegisterByPass(ctx, request)
	if err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
		return
	}

	messages.SendSuccessResponse(ctx, responses.SuccessResponse{
		SuccessCode: "SUCCESS-AUTH-0001",
	})
}

func (handler *AuthHttp) LoginByPass(ctx *gin.Context) {
	request := &requests.LoginRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	if formErrors, err := validators.ValidateStruct(ctx, request); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error:      err,
			FormErrors: formErrors,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	data, err := handler.authUc.LoginByPass(ctx, request)
	if err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
		return
	}

	messages.SendSuccessResponse(ctx, responses.SuccessResponse{
		SuccessCode: "SUCCESS-AUTH-0002",
		Data:        data,
	})
}
