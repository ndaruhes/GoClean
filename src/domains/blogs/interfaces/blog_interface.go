package interfaces

import (
	"context"
	"go-clean/src/domains/blogs/entities"
	"go-clean/src/models/requests"
	"go-clean/src/models/responses"
)

type BlogUseCase interface {
	// BLOG USECASE
	GetPublicBlogList(ctx context.Context, request *requests.BlogListRequest) ([]responses.PublicBlogListsResponse, int64, error)
	SearchBlog(ctx context.Context, request *requests.SearchBlogRequest) ([]responses.PublicBlogListsResponse, int64, error)
	GetBlogDetail(ctx context.Context, id string) (*responses.BlogDetailResponse, error)
	CreateBlog(ctx context.Context, request *requests.UpsertBlogRequest, file []byte, fileName string) error
	UpdateDraftedBlog(ctx context.Context, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error
	PublishBlog(ctx context.Context, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error
	UpdatePublishedBlog(ctx context.Context, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error
	UpdateSlug(ctx context.Context, blogID string, request *requests.UpdateSlugRequest) error
	UpdateBlogToDraft(ctx context.Context, blogID string) error
	DeleteBlog(ctx context.Context, blogID string) error
}

type BlogRepository interface {
	// BLOG REPOSITORY
	GetPublicBlogList(ctx context.Context, request *requests.BlogListRequest) ([]entities.Blog, int64, error)
	SearchBlog(ctx context.Context, request *requests.SearchBlogRequest) ([]entities.Blog, int64, error)
	GetBlogDetail(ctx context.Context, id string) (*entities.Blog, error)
	FindBlogById(ctx context.Context, id string) (*entities.Blog, error)
	FindBlogBySlug(ctx context.Context, slug string) (*entities.Blog, error)
	CreateBlog(ctx context.Context, blog *entities.Blog) (*entities.Blog, error)
	UpdateBlog(ctx context.Context, blogID string, blogStatusCheck string, blog *entities.Blog) error
	DeleteBlog(ctx context.Context, blogID string) error

	// BLOG CATEGORY REPOSITORY
	CreateBlogCategory(ctx context.Context, blogCategory []entities.BlogCategory) error
	UpdateBlogCategory(ctx context.Context, blogID string, blogCategory []entities.BlogCategory) error
}
