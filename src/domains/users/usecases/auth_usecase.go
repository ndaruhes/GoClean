package usecases

import (
	"go-clean/domains/users"
	"go-clean/domains/users/entities"
	"go-clean/models/requests"
	"go-clean/models/responses"
	errors2 "go-clean/shared/errors"
	"go-clean/shared/helpers"
	"go-clean/shared/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthUseCase struct {
	authRepo users.AuthRepository
}

func NewAuthUseCase(authRepo users.AuthRepository) *AuthUseCase {
	return &AuthUseCase{
		authRepo: authRepo,
	}
}

func (uc *AuthUseCase) RegisterWithEmailPassword(ctx *gin.Context, request *requests.RegisterWithEmailPasswordRequest) error {
	passConfig := &requests.PasswordConfig{
		Time:    1,
		Memory:  64 * 1024,
		Threads: 4,
		KeyLen:  32,
	}
	hashPassword, err := helpers.GeneratePassword(passConfig, request.Password)
	if err != nil {
		return err
	}

	newUser := &entities.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashPassword,
		GoogleID: nil,
	}

	if err := validators.ValidateStruct(ctx, newUser); err != nil {
		return &errors2.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusBadRequest,
		}
	}

	if err := uc.authRepo.RegisterWithEmailPassword(ctx, newUser); err != nil {
		return &errors2.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (uc *AuthUseCase) Login(ctx *gin.Context, request *requests.LoginRequest) (*responses.LoginResponse, error) {
	//TODO implement me
	panic("implement me")
}
