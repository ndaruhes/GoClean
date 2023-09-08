package repositories

import (
	"errors"
	"go-clean/src/domains/blogs/entities"
	"go-clean/src/models/messages"
	"go-clean/src/models/responses"
	"go-clean/src/shared/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
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
func (repo *BlogRepository) FindBlogById(ctx *fiber.Ctx, id string) (*entities.Blog, error) {
	fiberCtx := ctx.Context()
	var blog *entities.Blog
	if err := utils.GetDb(ctx, repo.db).WithContext(fiberCtx).Where("id = ?", id).First(&blog).Error; err != nil {
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

func (repo *BlogRepository) FindBlogBySlug(ctx *fiber.Ctx, slug string) (*entities.Blog, error) {
	fiberCtx := ctx.Context()
	var blog *entities.Blog
	if err := utils.GetDb(ctx, repo.db).WithContext(fiberCtx).Where("slug = ?", slug).First(&blog).Error; err != nil {
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

func (repo *BlogRepository) GetPublicBlogList(ctx *fiber.Ctx) (*responses.PublicBlogListsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *BlogRepository) CreateBlog(ctx *fiber.Ctx, blog *entities.Blog) (*entities.Blog, error) {
	fiberCtx := ctx.Context()
	err := utils.GetDb(ctx, repo.db).WithContext(fiberCtx).Create(&blog).Error
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (repo *BlogRepository) UpdateBlog(ctx *fiber.Ctx, blogID string, blogStatusCheck string, blog *entities.Blog) error {
	fiberCtx := ctx.Context()
	return utils.GetDb(ctx, repo.db).WithContext(fiberCtx).Model(&entities.Blog{}).Where("id = ?", blogID).Where("status = ?", blogStatusCheck).Updates(&blog).Error
}

func (repo *BlogRepository) DeleteBlog(ctx *fiber.Ctx, blogID string) error {
	fiberCtx := ctx.Context()
	blog := &entities.Blog{ID: blogID}
	return utils.GetDb(ctx, repo.db).WithContext(fiberCtx).Delete(blog).Error
}

// BLOG CATEGORY REPOSITORY
func (repo *BlogRepository) CreateBlogCategory(ctx *fiber.Ctx, blogCategory []entities.BlogCategory) error {
	fiberCtx := ctx.Context()
	if len(blogCategory) > 0 {
		return utils.GetDb(ctx, repo.db).WithContext(fiberCtx).Create(&blogCategory).Error
	}
	return nil
}

func (repo *BlogRepository) UpdateBlogCategory(ctx *fiber.Ctx, blogID string, blogCategory []entities.BlogCategory) error {
	fiberCtx := ctx.Context()
	if err := utils.GetDb(ctx, repo.db).WithContext(fiberCtx).Model(&entities.BlogCategory{}).Where("blog_id = ?", blogID).Delete(&entities.BlogCategory{}).Error; err != nil {
		return err
	}

	if len(blogCategory) > 0 {
		if err := utils.GetDb(ctx, repo.db).WithContext(fiberCtx).Create(&blogCategory).Error; err != nil {
			return err
		}
	}

	return nil
}
