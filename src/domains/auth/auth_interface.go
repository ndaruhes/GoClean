package auth

import (
	"go-clean/domains/users/entities"
	"go-clean/models/requests"
	"go-clean/models/responses"

	"github.com/gin-gonic/gin"
)

type AuthUseCase interface {
	LoginByPass(ctx *gin.Context, request *requests.LoginRequest) (*responses.LoginResponse, error)
	RegisterByPass(ctx *gin.Context, request *requests.RegisterWithEmailPasswordRequest) error
}

type AuthRepository interface {
	RegisterByPass(ctx *gin.Context, user *entities.User) error
	FindUserByOAuthTokenId(ctx *gin.Context, tokenId string) (*entities.User, error)
	FindByEmail(ctx *gin.Context, email string) (*entities.User, error)
	GenerateTokenUser(ctx *gin.Context, user *entities.User) (*entities.OAuthToken, error)
}
