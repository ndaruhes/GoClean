package repositories

import (
	"go-clean/src/domains/users/entities"
	"go-clean/src/models/messages"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (repo AuthRepository) RegisterByPass(ctx *fiber.Ctx, user *entities.User) error {
	fiberCtx := ctx.Context()
	return repo.db.WithContext(fiberCtx).Create(&user).Error
}

func (repo AuthRepository) FindUserByOAuthTokenId(ctx *fiber.Ctx, tokenId string) (*entities.User, error) {
	//TODO implement me
	panic("implement me")
}

func (repo AuthRepository) FindByEmail(ctx *fiber.Ctx, email string) (*entities.User, error) {
	fiberCtx := ctx.Context()
	var user *entities.User
	if err := repo.db.WithContext(fiberCtx).Model(&entities.User{}).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (repo AuthRepository) GenerateTokenUser(ctx *fiber.Ctx, user *entities.User) (*entities.OAuthToken, error) {
	fiberCtx := ctx.Context()
	jwtToken := &entities.OAuthToken{
		ExpiresAt: time.Now().AddDate(0, 1, 0),
		UserID:    user.ID,
	}

	err := repo.db.WithContext(fiberCtx).Save(&jwtToken).Error
	if messages.HasError(err) {
		return nil, &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return jwtToken, nil
}
