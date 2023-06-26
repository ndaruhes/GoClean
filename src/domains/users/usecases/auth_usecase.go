package usecases

import (
	"errors"
	"go-clean/domains/users"
	"go-clean/domains/users/entities"
	errors2 "go-clean/models/messages"
	"go-clean/models/requests"
	"go-clean/models/responses"
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

func (uc *AuthUseCase) RegisterByPass(ctx *gin.Context, request *requests.RegisterWithEmailPasswordRequest) error {
	hashPassword, err := helpers.GeneratePassword(request.Password)
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

	if err := uc.authRepo.RegisterByPass(ctx, newUser); err != nil {
		return &errors2.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (uc *AuthUseCase) LoginByPass(ctx *gin.Context, request *requests.LoginRequest) (*responses.LoginResponse, error) {
	user, err := uc.authRepo.FindByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	passwordPassed, err := helpers.ComparePassword(request.Password, user.Password)
	if err != nil {
		return nil, err
	}

	if passwordPassed == false {
		return nil, &errors2.ErrorWrapper{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("Wrong Credentials"),
		}
	}

	token, err := helpers.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	_, err = uc.authRepo.GenerateTokenUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &responses.LoginResponse{
		TokenID: token,
	}, nil
}
