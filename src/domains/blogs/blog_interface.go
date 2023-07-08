package blogs

import (
	"github.com/gin-gonic/gin"
	"go-clean/domains/blogs/entities"
	"go-clean/models/requests"
	"go-clean/models/responses"
)

type BlogUseCase interface {
	GetPublicBlogList(ctx *gin.Context) (*responses.PublicBlogListsResponse, error)
	CreateBlog(ctx *gin.Context, request *requests.UpsertBlogRequest, file []byte, fileName string) error
	UpdateBlog(ctx *gin.Context, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error
}

type BlogRepository interface {
	GetPublicBlogList(ctx *gin.Context) (*responses.PublicBlogListsResponse, error)
	CreateBlog(ctx *gin.Context, blog *entities.Blog) error
	FindBlogBySlug(ctx *gin.Context, slug string) (*entities.Blog, error)
}
