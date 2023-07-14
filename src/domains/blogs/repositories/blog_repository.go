package repositories

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-clean/domains/blogs/entities"
	"go-clean/models/messages"
	"go-clean/models/responses"
	"gorm.io/gorm"
	"net/http"
)

type BlogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) *BlogRepository {
	return &BlogRepository{
		db: db,
	}
}

func (repo *BlogRepository) BeginTransaction(ctx *gin.Context) *gorm.DB {
	db := repo.db.Begin()
	ctx.Set("tx", db)
	return db
}

func (repo *BlogRepository) Commit(ctx *gin.Context) {
	db, _ := ctx.Get("tx")
	if tx, ok := db.(*gorm.DB); ok {
		tx.Commit()
	}
}

func (repo *BlogRepository) FindBlogById(ctx *gin.Context, id string) (*entities.Blog, error) {
	var blog *entities.Blog
	if err := repo.db.WithContext(ctx).Where("id = ?", id).First(&blog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &messages.ErrorWrapper{
				Context:    ctx,
				Err:        err,
				ErrorCode:  "ERROR-404001",
				StatusCode: http.StatusNotFound,
			}
		}
	}

	return blog, nil
}

func (repo *BlogRepository) FindBlogBySlug(ctx *gin.Context, slug string) (*entities.Blog, error) {
	var blog *entities.Blog
	if err := repo.db.WithContext(ctx).Where("slug = ?", slug).First(&blog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &messages.ErrorWrapper{
				Context:    ctx,
				Err:        err,
				ErrorCode:  "ERROR-404001",
				StatusCode: http.StatusNotFound,
			}
		}
	}

	return blog, nil
}

func (repo *BlogRepository) GetPublicBlogList(ctx *gin.Context) (*responses.PublicBlogListsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *BlogRepository) CreateBlog(ctx *gin.Context, blog *entities.Blog) (*entities.Blog, error) {
	err := repo.db.WithContext(ctx).Create(&blog).Error
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (repo *BlogRepository) CreateBlogCategory(ctx *gin.Context, blogCategory []entities.BlogCategory) error {
	return repo.db.WithContext(ctx).Create(&blogCategory).Error
}

func (repo *BlogRepository) UpdateBlog(ctx *gin.Context, blogID string, blogStatusCheck string, blog *entities.Blog) error {
	return repo.db.WithContext(ctx).Model(&entities.Blog{}).Where("id = ?", blogID).Where("status = ?", blogStatusCheck).Updates(&blog).Error
}

func (repo *BlogRepository) DeleteBlog(ctx *gin.Context, blogID string) error {
	blog := &entities.Blog{ID: blogID}
	return repo.db.WithContext(ctx).Delete(blog).Error
}
