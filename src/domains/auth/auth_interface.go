package auth

import (
	"go-clean/src/domains/users/entities"
	"go-clean/src/models/requests"
	"go-clean/src/models/responses"

	"github.com/gofiber/fiber/v2"
)

type AuthUseCase interface {
	LoginByPass(ctx *fiber.Ctx, request *requests.LoginRequest) (*responses.LoginResponse, error)
	RegisterByPass(ctx *fiber.Ctx, request *requests.RegisterWithEmailPasswordRequest) error
}

type AuthRepository interface {
	RegisterByPass(ctx *fiber.Ctx, user *entities.User) error
	FindUserByOAuthTokenId(ctx *fiber.Ctx, tokenId string) (*entities.User, error)
	FindByEmail(ctx *fiber.Ctx, email string) (*entities.User, error)
	GenerateTokenUser(ctx *fiber.Ctx, user *entities.User) (*entities.OAuthToken, error)
}
