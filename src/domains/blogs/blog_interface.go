package blogs

import (
	"github.com/gin-gonic/gin"
	"go-clean/domains/blogs/entities"
	"go-clean/models/requests"
	"go-clean/models/responses"
	"gorm.io/gorm"
)

type BlogUseCase interface {
	GetPublicBlogList(ctx *gin.Context) (*responses.PublicBlogListsResponse, error)
	CreateBlog(ctx *gin.Context, request *requests.UpsertBlogRequest, file []byte, fileName string) error
	AdjustBlog(ctx *gin.Context, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error
	PublishBlog(ctx *gin.Context, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error
	UpdateBlog(ctx *gin.Context, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error
	UpdateSlug(ctx *gin.Context, blogID string, request *requests.UpsertBlogRequest) error
	UpdateBlogToDraft(ctx *gin.Context, blogID string) error
	DeleteBlog(ctx *gin.Context, blogID string) error
}

type BlogRepository interface {
	BeginTransaction(ctx *gin.Context) *gorm.DB
	Commit(ctx *gin.Context)
	FindBlogById(ctx *gin.Context, id string) (*entities.Blog, error)
	FindBlogBySlug(ctx *gin.Context, slug string) (*entities.Blog, error)
	GetPublicBlogList(ctx *gin.Context) (*responses.PublicBlogListsResponse, error)
	CreateBlog(ctx *gin.Context, blog *entities.Blog) (*entities.Blog, error)
	CreateBlogCategory(ctx *gin.Context, blogCategory []entities.BlogCategory) error
	UpdateBlog(ctx *gin.Context, blogID string, blogStatusCheck string, blog *entities.Blog) error
	DeleteBlog(ctx *gin.Context, blogID string) error
}
