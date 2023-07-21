package usecases

import (
	"errors"
	"go-clean/domains/auth"
	"go-clean/domains/users/entities"
	"go-clean/models/messages"
	"go-clean/models/requests"
	"go-clean/models/responses"
	"go-clean/shared/helpers"
	"go-clean/shared/validators"
	"gorm.io/gorm"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthUseCase struct {
	authRepo auth.AuthRepository
}

func NewAuthUseCase(authRepo auth.AuthRepository) *AuthUseCase {
	return &AuthUseCase{
		authRepo: authRepo,
	}
}

func (uc *AuthUseCase) RegisterByPass(ctx *gin.Context, request *requests.RegisterWithEmailPasswordRequest) error {
	user, err := uc.authRepo.FindByEmail(ctx, request.Email)
	if user != nil {
		return &messages.ErrorWrapper{
			Context:    ctx,
			ErrorCode:  "ERROR-401001",
			StatusCode: http.StatusUnauthorized,
		}
	}

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

	if _, err := validators.ValidateStruct(ctx, newUser); err != nil {
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusBadRequest,
		}
	}

	if err := uc.authRepo.RegisterByPass(ctx, newUser); err != nil {
		return err
	}

	return nil
}

func (uc *AuthUseCase) LoginByPass(ctx *gin.Context, request *requests.LoginRequest) (*responses.LoginResponse, error) {
	user, err := uc.authRepo.FindByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &messages.ErrorWrapper{
				Context:    ctx,
				ErrorCode:  "ERROR-401002",
				StatusCode: http.StatusUnauthorized,
			}
		}
	}

	passwordPassed, err := helpers.ComparePassword(request.Password, user.Password)
	if err != nil {
		return nil, err
	}

	if passwordPassed == false {
		return nil, &messages.ErrorWrapper{
			Context:    ctx,
			ErrorCode:  "ERROR-401003",
			StatusCode: http.StatusUnauthorized,
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
		Token: token,
	}, nil
}
