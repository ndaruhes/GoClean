package blogs

import (
	"github.com/gin-gonic/gin"
	"go-clean/domains/blogs/entities"
	"go-clean/models/requests"
	"go-clean/models/responses"
)

type BlogUseCase interface {
	// BLOG USECASE
	GetPublicBlogList(ctx *gin.Context) (*responses.PublicBlogListsResponse, error)
	CreateBlog(ctx *gin.Context, request *requests.UpsertBlogRequest, file []byte, fileName string) error
	AdjustBlog(ctx *gin.Context, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error
	PublishBlog(ctx *gin.Context, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error
	UpdateBlog(ctx *gin.Context, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error
	UpdateSlug(ctx *gin.Context, blogID string, request *requests.UpdateSlugRequest) error
	UpdateBlogToDraft(ctx *gin.Context, blogID string) error
	DeleteBlog(ctx *gin.Context, blogID string) error
}

type BlogRepository interface {
	// BLOG REPOSITORY
	FindBlogById(ctx *gin.Context, id string) (*entities.Blog, error)
	FindBlogBySlug(ctx *gin.Context, slug string) (*entities.Blog, error)
	GetPublicBlogList(ctx *gin.Context) (*responses.PublicBlogListsResponse, error)
	CreateBlog(ctx *gin.Context, blog *entities.Blog) (*entities.Blog, error)
	UpdateBlog(ctx *gin.Context, blogID string, blogStatusCheck string, blog *entities.Blog) error
	DeleteBlog(ctx *gin.Context, blogID string) error

	// BLOG CATEGORY REPOSITORY
	CreateBlogCategory(ctx *gin.Context, blogCategory []entities.BlogCategory) error
	UpdateBlogCategory(ctx *gin.Context, blogID string, blogCategory []entities.BlogCategory) error
}
