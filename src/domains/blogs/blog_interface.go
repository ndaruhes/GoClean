package blogs

import (
	"go-clean/src/domains/blogs/entities"
	"go-clean/src/models/requests"
	"go-clean/src/models/responses"

	"github.com/gofiber/fiber/v2"
)

type BlogUseCase interface {
	// BLOG USECASE
	GetPublicBlogList(ctx *fiber.Ctx) (*responses.PublicBlogListsResponse, error)
	CreateBlog(ctx *fiber.Ctx, request *requests.UpsertBlogRequest, file []byte, fileName string) error
	AdjustBlog(ctx *fiber.Ctx, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error
	PublishBlog(ctx *fiber.Ctx, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error
	UpdateBlog(ctx *fiber.Ctx, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error
	UpdateSlug(ctx *fiber.Ctx, blogID string, request *requests.UpdateSlugRequest) error
	UpdateBlogToDraft(ctx *fiber.Ctx, blogID string) error
	DeleteBlog(ctx *fiber.Ctx, blogID string) error
}

type BlogRepository interface {
	// BLOG REPOSITORY
	FindBlogById(ctx *fiber.Ctx, id string) (*entities.Blog, error)
	FindBlogBySlug(ctx *fiber.Ctx, slug string) (*entities.Blog, error)
	GetPublicBlogList(ctx *fiber.Ctx) (*responses.PublicBlogListsResponse, error)
	CreateBlog(ctx *fiber.Ctx, blog *entities.Blog) (*entities.Blog, error)
	UpdateBlog(ctx *fiber.Ctx, blogID string, blogStatusCheck string, blog *entities.Blog) error
	DeleteBlog(ctx *fiber.Ctx, blogID string) error

	// BLOG CATEGORY REPOSITORY
	CreateBlogCategory(ctx *fiber.Ctx, blogCategory []entities.BlogCategory) error
	UpdateBlogCategory(ctx *fiber.Ctx, blogID string, blogCategory []entities.BlogCategory) error
}
