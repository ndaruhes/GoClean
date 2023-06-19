package repositories

import (
	"go-clean/domains/users/entities"

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

func (repo AuthRepository) RegisterWithEmailPassword(ctx *gin.Context, user *entities.User) error {
	return repo.db.WithContext(ctx).Create(&user).Error
}

func (repo AuthRepository) FindUserByOAuthTokenId(ctx *gin.Context, tokenId string) (*entities.User, error) {
	//TODO implement me
	panic("implement me")
}

func (repo AuthRepository) FindByEmail(email string) (*entities.User, error) {
	//TODO implement me
	panic("implement me")
}

func (repo AuthRepository) GenerateToken(ctx *gin.Context, user *entities.User) (*entities.OAuthToken, error) {
	//TODO implement me
	panic("implement me")
}
