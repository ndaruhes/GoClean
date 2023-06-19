package users

import (
	"go-clean/domains/users/entities"
	"go-clean/models/requests"
	"go-clean/models/responses"

	"github.com/gin-gonic/gin"
)

type AuthUseCase interface {
	Login(ctx *gin.Context, request *requests.LoginRequest) (*responses.LoginResponse, error)
	RegisterWithEmailPassword(ctx *gin.Context, request *requests.RegisterWithEmailPasswordRequest) error
}

type AuthRepository interface {
	RegisterWithEmailPassword(ctx *gin.Context, user *entities.User) error
	FindUserByOAuthTokenId(ctx *gin.Context, tokenId string) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	GenerateToken(ctx *gin.Context, user *entities.User) (*entities.OAuthToken, error)
}
