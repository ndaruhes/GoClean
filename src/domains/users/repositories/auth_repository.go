package repositories

import (
	"go-clean/domains/users/entities"
	errors2 "go-clean/models/messages"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func (repo AuthRepository) RegisterByPass(ctx *gin.Context, user *entities.User) error {
	return repo.db.WithContext(ctx).Create(&user).Error
}

func (repo AuthRepository) FindUserByOAuthTokenId(ctx *gin.Context, tokenId string) (*entities.User, error) {
	//TODO implement me
	panic("implement me")
}

func (repo AuthRepository) FindByEmail(ctx *gin.Context, email string) (*entities.User, error) {
	var user *entities.User
	if err := repo.db.WithContext(ctx).Model(&entities.User{}).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (repo AuthRepository) GenerateTokenUser(ctx *gin.Context, user *entities.User) (*entities.OAuthToken, error) {
	jwtToken := &entities.OAuthToken{
		ExpiresAt: time.Now().AddDate(0, 1, 0),
		UserID:    user.ID,
	}

	err := repo.db.WithContext(ctx).Save(&jwtToken).Error
	if errors2.HasError(err) {
		return nil, &errors2.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return jwtToken, nil
}
