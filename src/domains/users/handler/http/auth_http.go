package http

import (
	"go-clean/configs/database"
	"go-clean/domains/users"
	authRepository "go-clean/domains/users/repositories"
	authUseCase "go-clean/domains/users/usecases"
	"go-clean/models/requests"
	"go-clean/models/responses"
	"go-clean/shared/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHttp struct {
	authUc users.AuthUseCase
}

func NewAuthHttp(route *gin.Engine) *AuthHttp {
	handler := &AuthHttp{
		authUc: authUseCase.NewAuthUseCase(
			authRepository.NewAuthRepository(database.ConnectDatabase()),
		),
	}

	auth := route.Group("auth")
	{
		auth.POST("/register", handler.RegisterWithEmailPassword)
	}

	return handler
}

func (handler *AuthHttp) RegisterWithEmailPassword(ctx *gin.Context) {
	request := &requests.RegisterWithEmailPasswordRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, responses.BasicResponse{
			Message: "Error Validation",
			Error:   err,
		})
	}

	if err := validators.ValidateStruct(ctx, request); err != nil {
		ctx.JSON(http.StatusBadRequest, responses.BasicResponse{
			Message: "Error Validation",
			Error:   err,
		})
	}

	err := handler.authUc.RegisterWithEmailPassword(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responses.BasicResponse{
			Error: err,
		})
	}
}
