package http

import (
	"go-clean/configs/database"
	"go-clean/domains/blogs"
	blogRepository "go-clean/domains/blogs/repositories"
	blogUseCase "go-clean/domains/blogs/usecases"
	"go-clean/middlewares"
	"go-clean/models/messages"
	"go-clean/models/requests"
	"go-clean/models/responses"
	"go-clean/shared/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BlogHttp struct {
	blogUc blogs.BlogUseCase
}

func NewBlogHttp(route *gin.Engine) *BlogHttp {
	handler := &BlogHttp{
		blogUc: blogUseCase.NewBlogUseCase(blogRepository.NewBlogRepository(database.ConnectDatabase())),
	}

	blog := route.Group("blog")
	{
		blog.Use(middlewares.Authenticated())
		blog.POST("", handler.CreateBlog)
	}

	return handler
}

func (handler *BlogHttp) CreateBlog(ctx *gin.Context) {
	request := &requests.UpsertBlogRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	if formErrors, err := validators.ValidateStruct(ctx, request); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error:      err,
			FormErrors: formErrors,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	err := handler.blogUc.CreateBlog(ctx, request)
	if err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
		return
	}

	messages.SendSuccessResponse(ctx, responses.SuccessResponse{
		SuccessCode: "SUCCESS-BLOG-0001",
	})
}
