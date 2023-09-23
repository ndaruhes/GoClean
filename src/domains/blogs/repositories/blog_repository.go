package repositories

import (
	"context"
	"errors"
	"go-clean/src/domains/blogs/entities"
	"go-clean/src/models/messages"
	"go-clean/src/shared/utils"
	"net/http"

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

// BLOG REPOSITORY
func (repo *BlogRepository) FindBlogById(ctx context.Context, id string) (*entities.Blog, error) {
	var blog *entities.Blog
	if err := utils.GetDb(ctx, repo.db).WithContext(ctx).Where("id = ?", id).First(&blog).Error; err != nil {
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

func (repo *BlogRepository) FindBlogBySlug(ctx context.Context, slug string) (*entities.Blog, error) {
	var blog *entities.Blog
	if err := utils.GetDb(ctx, repo.db).WithContext(ctx).Where("slug = ?", slug).First(&blog).Error; err != nil {
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

func (repo *BlogRepository) GetPublicBlogList(ctx context.Context) ([]entities.Blog, error) {
	var data []entities.Blog
	err := utils.GetDb(ctx, repo.db).WithContext(ctx).
		Model(&entities.Blog{}).
		Preload("User").
		Where("status = ?", "Published").
		Where("published_at is not null").
		Find(&data).Error

	return data, err
}

func (repo *BlogRepository) CreateBlog(ctx context.Context, blog *entities.Blog) (*entities.Blog, error) {
	err := utils.GetDb(ctx, repo.db).WithContext(ctx).Create(&blog).Error
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (repo *BlogRepository) UpdateBlog(ctx context.Context, blogID string, blogStatusCheck string, blog *entities.Blog) error {
	return utils.GetDb(ctx, repo.db).WithContext(ctx).Model(&entities.Blog{}).Where("id = ?", blogID).Where("status = ?", blogStatusCheck).Updates(&blog).Error
}

func (repo *BlogRepository) DeleteBlog(ctx context.Context, blogID string) error {
	blog := &entities.Blog{ID: blogID}
	return utils.GetDb(ctx, repo.db).WithContext(ctx).Delete(blog).Error
}

// BLOG CATEGORY REPOSITORY
func (repo *BlogRepository) CreateBlogCategory(ctx context.Context, blogCategory []entities.BlogCategory) error {
	if len(blogCategory) > 0 {
		return utils.GetDb(ctx, repo.db).WithContext(ctx).Create(&blogCategory).Error
	}
	return nil
}

func (repo *BlogRepository) UpdateBlogCategory(ctx context.Context, blogID string, blogCategory []entities.BlogCategory) error {
	if err := utils.GetDb(ctx, repo.db).WithContext(ctx).Model(&entities.BlogCategory{}).Where("blog_id = ?", blogID).Delete(&entities.BlogCategory{}).Error; err != nil {
		return err
	}

	if len(blogCategory) > 0 {
		if err := utils.GetDb(ctx, repo.db).WithContext(ctx).Create(&blogCategory).Error; err != nil {
			return err
		}
	}

	return nil
}
