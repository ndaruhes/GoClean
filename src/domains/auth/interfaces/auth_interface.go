package interfaces

import (
	"context"
	"go-clean/src/domains/users/entities"
	"go-clean/src/models/requests"
	"go-clean/src/models/responses"
)

type AuthUseCase interface {
	LoginByPass(ctx context.Context, request *requests.LoginRequest) (*responses.LoginResponse, error)
	RegisterByPass(ctx context.Context, request *requests.RegisterWithEmailPasswordRequest) error
}

type AuthRepository interface {
	RegisterByPass(ctx context.Context, user *entities.User) error
	FindUserByOAuthTokenId(ctx context.Context, tokenId string) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	GenerateTokenUser(ctx context.Context, user *entities.User) (*entities.OAuthToken, error)
}
