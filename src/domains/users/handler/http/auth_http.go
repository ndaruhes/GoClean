package http

import (
	"go-clean/configs/database"
	"go-clean/domains/users"
	authRepository "go-clean/domains/users/repositories"
	authUseCase "go-clean/domains/users/usecases"
	"go-clean/models/messages"
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
		auth.POST("/login", handler.LoginByPass)
	}

	return handler
}

func (handler *AuthHttp) RegisterWithEmailPassword(ctx *gin.Context) {
	request := &requests.RegisterWithEmailPasswordRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		messages.SendErrorResponse(ctx, responses.BasicResponse{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	if err := validators.ValidateStruct(ctx, request); err != nil {
		messages.SendErrorResponse(ctx, responses.BasicResponse{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	err := handler.authUc.RegisterByPass(ctx, request)
	if err != nil {
		messages.SendErrorResponse(ctx, responses.BasicResponse{
			Error:       err,
			MessageCode: "ERROR-400003",
			StatusCode:  http.StatusBadRequest,
		})
		return
	}

	messages.SendSuccessResponse(ctx, responses.BasicResponse{
		MessageCode: "SUCCESS-0001",
	})
}

func (handler *AuthHttp) LoginByPass(ctx *gin.Context) {
	request := &requests.LoginRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, responses.BasicResponse{
			Error: err,
		})
	}

	if err := validators.ValidateStruct(ctx, request); err != nil {
		ctx.JSON(http.StatusBadRequest, responses.BasicResponse{
			Error: err,
		})
	}

	res, err := handler.authUc.LoginByPass(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responses.BasicResponse{
			Error: err,
		})
	}

	ctx.JSON(http.StatusOK, responses.LoginResponse{
		TokenID: res.TokenID,
	})
}
