package http

import (
	"go-clean/src/domains/blogs"
	blogRepository "go-clean/src/domains/blogs/repositories"
	blogUseCase "go-clean/src/domains/blogs/usecases"
	"go-clean/src/middlewares"
	"go-clean/src/models/messages"
	"go-clean/src/models/requests"
	"go-clean/src/models/responses"
	"go-clean/src/setup/database"
	"go-clean/src/shared/utils"
	"go-clean/src/shared/validators"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BlogHttp struct {
	blogUc blogs.BlogUseCase
}

func NewBlogHttp(route *gin.Engine) *BlogHttp {
	db := database.ConnectDatabase()
	blogRepo := blogRepository.NewBlogRepository(db)
	blogUc := blogUseCase.NewBlogUseCase(blogRepo, db)

	handler := &BlogHttp{blogUc: blogUc}
	setRoutes(route, handler)

	return handler
}

func setRoutes(route *gin.Engine, handler *BlogHttp) {
	blog := route.Group("blog")
	{
		blog.Use(middlewares.Authenticated())
		blog.POST("", handler.CreateBlog)
		blog.PUT("/:id/edit", handler.AdjustBlog)
		blog.PUT("/:id/publish", handler.PublishBlog)
		blog.PUT("/:id", handler.UpdateBlog)
		blog.DELETE("/:id", handler.DeleteBlog)
		blog.PUT("/:id/slug", handler.UpdateSlug)
		blog.PUT("/:id/draft", handler.UpdateToDraft)
	}
}

func (handler *BlogHttp) CreateBlog(ctx *gin.Context) {
	request := &requests.UpsertBlogRequest{
		Title:   ctx.PostForm("title"),
		Content: ctx.PostForm("content"),
	}

	blogCategoryIds := handleBlogCategories(ctx)
	request.BlogCategoryIds = blogCategoryIds

	file, fileName := handleSingleFile(ctx)
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

func (handler *BlogHttp) AdjustBlog(ctx *gin.Context) {
	request := &requests.UpsertBlogRequest{
		Title:   ctx.PostForm("title"),
		Content: ctx.PostForm("content"),
	}

	blogCategoryIds := handleBlogCategories(ctx)
	request.BlogCategoryIds = blogCategoryIds

	file, fileName := handleSingleFile(ctx)
	if err := handler.blogUc.AdjustBlog(ctx, ctx.Param("id"), request, file, fileName); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
		return
	}

	messages.SendSuccessResponse(ctx, responses.SuccessResponse{
		SuccessCode: "SUCCESS-BLOG-0002",
	})
}

func (handler *BlogHttp) PublishBlog(ctx *gin.Context) {
	request := &requests.UpsertBlogRequest{
		Title:   ctx.PostForm("title"),
		Content: ctx.PostForm("content"),
	}

	blogCategoryIds := handleBlogCategories(ctx)
	request.BlogCategoryIds = blogCategoryIds

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

	file, fileName := handleSingleFile(ctx)
	if err := handler.blogUc.PublishBlog(ctx, ctx.Param("id"), request, file, fileName); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
		return
	}

	messages.SendSuccessResponse(ctx, responses.SuccessResponse{
		SuccessCode: "SUCCESS-BLOG-0004",
	})
}

func (handler *BlogHttp) UpdateBlog(ctx *gin.Context) {
	request := &requests.UpsertBlogRequest{
		Title:   ctx.PostForm("title"),
		Content: ctx.PostForm("content"),
	}

	blogCategoryIds := handleBlogCategories(ctx)
	request.BlogCategoryIds = blogCategoryIds

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

	file, fileName := handleSingleFile(ctx)
	if err := handler.blogUc.UpdateBlog(ctx, ctx.Param("id"), request, file, fileName); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
		return
	}

	messages.SendSuccessResponse(ctx, responses.SuccessResponse{
		SuccessCode: "SUCCESS-BLOG-0002",
	})
}

func (handler *BlogHttp) DeleteBlog(ctx *gin.Context) {
	if err := handler.blogUc.DeleteBlog(ctx, ctx.Param("id")); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
		return
	}

	messages.SendSuccessResponse(ctx, responses.SuccessResponse{
		SuccessCode: "SUCCESS-BLOG-0003",
	})
}

func (handler *BlogHttp) UpdateSlug(ctx *gin.Context) {
	request := &requests.UpdateSlugRequest{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
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

	if err := handler.blogUc.UpdateSlug(ctx, ctx.Param("id"), request); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
		return
	}

	messages.SendSuccessResponse(ctx, responses.SuccessResponse{
		SuccessCode: "SUCCESS-BLOG-0005",
	})
}

func (handler *BlogHttp) UpdateToDraft(ctx *gin.Context) {
	if err := handler.blogUc.UpdateBlogToDraft(ctx, ctx.Param("id")); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
		return
	}

	messages.SendSuccessResponse(ctx, responses.SuccessResponse{
		SuccessCode: "SUCCESS-BLOG-0006",
	})
}

func handleBlogCategories(ctx *gin.Context) []int {
	blogCategoryIds, _ := ctx.GetPostFormArray("blogCategoryIds[]")
	convertedCategoryIds := make([]int, len(blogCategoryIds))
	for i, id := range blogCategoryIds {
		convertedId, err := strconv.Atoi(id)
		if err != nil {
			messages.SendErrorResponse(ctx, responses.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Error:      err,
			})
			return nil
		}
		convertedCategoryIds[i] = convertedId
	}

	return convertedCategoryIds
}

func handleSingleFile(ctx *gin.Context) ([]byte, string) {
	var (
		file     []byte
		fileName string
	)
	header, err := ctx.FormFile("cover")
	if header != nil {
		if messages.HasError(err) {
			messages.SendErrorResponse(ctx, responses.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Error:      err,
			})
			return nil, ""
		}

		err = validators.ValidateImage(header)
		if messages.HasError(err) {
			messages.SendErrorResponse(ctx, responses.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Error:      err,
			})
			return nil, ""
		}

		file, err = utils.MultipartFileHeaderToByte(header)
		if messages.HasError(err) {
			messages.SendErrorResponse(ctx, responses.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Error:      err,
			})
			return nil, ""
		}

		fileName = utils.GenerateFileName(header)
	}

	return file, fileName
}
