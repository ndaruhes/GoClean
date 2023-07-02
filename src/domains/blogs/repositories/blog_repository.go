package repositories

import (
	"github.com/gin-gonic/gin"
	"go-clean/domains/blogs/entities"
	"go-clean/models/responses"
	"gorm.io/gorm"
)

type BlogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) *BlogRepository {
	return &BlogRepository{
		db: db,
	}
}

func (repo BlogRepository) GetPublicBlogList(ctx *gin.Context) (*responses.PublicBlogListsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (repo BlogRepository) CreateBlog(ctx *gin.Context, blog *entities.Blog) error {
	return repo.db.Create(&blog).Error
}

func (repo BlogRepository) FindBlogBySlug(ctx *gin.Context, slug string) (*entities.Blog, error) {
	var blog *entities.Blog
	if err := repo.db.Where("slug = ?", slug).First(&blog).Error; err != nil {
		return nil, err
	}

	return blog, nil
}
