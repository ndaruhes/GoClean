package usecases

import (
	"context"
	"errors"
	"go-clean/src/domains/auth/interfaces"
	"go-clean/src/domains/users/entities"
	"go-clean/src/models/messages"
	"go-clean/src/models/requests"
	"go-clean/src/models/responses"
	"go-clean/src/shared/helpers"
	"go-clean/src/shared/validators"
	"net/http"

	"gorm.io/gorm"
)

type AuthUseCase struct {
	authRepo interfaces.AuthRepository
}

func NewAuthUseCase(authRepo interfaces.AuthRepository) *AuthUseCase {
	return &AuthUseCase{
		authRepo: authRepo,
	}
}

func (uc *AuthUseCase) RegisterByPass(ctx context.Context, request *requests.RegisterWithEmailPasswordRequest) error {
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

func (uc *AuthUseCase) LoginByPass(ctx context.Context, request *requests.LoginRequest) (*responses.LoginResponse, error) {
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
		Role:  user.Role,
		Token: token,
	}, nil
}
