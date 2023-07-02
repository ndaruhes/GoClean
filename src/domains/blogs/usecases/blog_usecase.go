package usecases

import (
	"github.com/gin-gonic/gin"
	"go-clean/domains/blogs"
	"go-clean/domains/blogs/entities"
	"go-clean/models/messages"
	"go-clean/models/requests"
	"go-clean/models/responses"
	"go-clean/shared/utils"
	"go-clean/shared/validators"
	"net/http"
)

type BlogUseCase struct {
	blogRepo blogs.BlogRepository
}

func NewBlogUseCase(blogRepo blogs.BlogRepository) *BlogUseCase {
	return &BlogUseCase{
		blogRepo: blogRepo,
	}
}

func (uc BlogUseCase) GetPublicBlogList(ctx *gin.Context) (*responses.PublicBlogListsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (uc BlogUseCase) CreateBlog(ctx *gin.Context, request *requests.UpsertBlogRequest) error {
	user := ctx.Value("member").(*responses.TokenDecoded)
	newBlog := &entities.Blog{
		Title:   request.Title,
		Cover:   request.Cover,
		Slug:    utils.GenerateSlug(request.Title),
		Content: request.Content,
		UserID:  user.ID,
	}

	if _, err := validators.ValidateStruct(ctx, newBlog); err != nil {
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusBadRequest,
		}
	}

	if err := uc.blogRepo.CreateBlog(ctx, newBlog); err != nil {
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}
