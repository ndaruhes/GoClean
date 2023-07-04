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
	"go-clean/shared/utils"
	"go-clean/shared/validators"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
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
	request := &requests.UpsertBlogRequest{
		Title:   ctx.PostForm("title"),
		Content: ctx.PostForm("content"),
	}
	if err := ctx.ShouldBind(&request); err != nil {
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

	header, err := ctx.FormFile("cover")
	//var oke, header, err = ctx.Request.FormFile("cover")
	if messages.HasError(err) {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err,
		})
		return
	}

	//err = validators.ValidateImage(header, validations.ImageValidation{
	//	MaxSize:   2,
	//	MinWidth:  300,
	//	MaxWidth:  640,
	//	MinHeight: 300,
	//	MaxHeight: 640,
	//	Format:    []string{"jpeg", "png", "jpg"},
	//})
	//if messages.HasError(err) {
	//	messages.SendErrorResponse(ctx, responses.ErrorResponse{
	//		StatusCode: http.StatusBadRequest,
	//		Error:      err,
	//	})
	//	return
	//}

	file, err := utils.MultipartFileHeaderToByte(header)
	if messages.HasError(err) {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err,
		})
		return
	}

	fileName := strings.ToUpper(xid.New().String()) + "-" + header.Filename
	if err := handler.blogUc.CreateBlog(ctx, request, file, fileName); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
		return
	}

	messages.SendSuccessResponse(ctx, responses.SuccessResponse{
		SuccessCode: "SUCCESS-BLOG-0001",
	})
}
