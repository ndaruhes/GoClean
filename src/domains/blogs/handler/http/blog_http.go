package http

import (
	"github.com/gofiber/fiber/v2"
	"go-clean/src/app/infrastructures"
	"go-clean/src/domains/blogs"
	blogRepository "go-clean/src/domains/blogs/repositories"
	blogUseCase "go-clean/src/domains/blogs/usecases"
	"go-clean/src/middlewares"
	"go-clean/src/models/messages"
	"go-clean/src/models/requests"
	"go-clean/src/models/responses"
	"go-clean/src/shared/utils"
	"go-clean/src/shared/validators"
	"net/http"
	"strconv"
	"strings"
)

type BlogHttp struct {
	blogUc blogs.BlogUseCase
}

func NewBlogHttp(route *fiber.App) *BlogHttp {
	db := infrastructures.ConnectDatabase()
	blogRepo := blogRepository.NewBlogRepository(db)
	blogUc := blogUseCase.NewBlogUseCase(blogRepo, db)

	handler := &BlogHttp{blogUc: blogUc}
	setRoutes(route, handler)

	return handler
}

func setRoutes(route *fiber.App, handler *BlogHttp) {
	blog := route.Group("blog")
	{
		blog.Use(middlewares.Authenticated())
		blog.Post("", handler.CreateBlog)
		blog.Put("/:id/edit", handler.AdjustBlog)
		blog.Put("/:id/publish", handler.PublishBlog)
		blog.Put("/:id", handler.UpdateBlog)
		blog.Delete("/:id", handler.DeleteBlog)
		blog.Put("/:id/slug", handler.UpdateSlug)
		blog.Put("/:id/draft", handler.UpdateToDraft)
	}
}

func (handler *BlogHttp) CreateBlog(ctx *fiber.Ctx) error {
	request := &requests.UpsertBlogRequest{
		Title:   ctx.FormValue("title"),
		Content: ctx.FormValue("content"),
	}

	blogCategoryIds := handleBlogCategories(ctx)
	request.BlogCategoryIds = blogCategoryIds

	file, fileName := handleSingleFile(ctx)
	if err := handler.blogUc.CreateBlog(ctx, request, file, fileName); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
	} else {
		messages.SendSuccessResponse(ctx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-BLOG-0001",
		})
	}

	return nil
}

func (handler *BlogHttp) AdjustBlog(ctx *fiber.Ctx) error {
	request := &requests.UpsertBlogRequest{
		Title:   ctx.FormValue("title"),
		Content: ctx.FormValue("content"),
	}

	blogCategoryIds := handleBlogCategories(ctx)
	request.BlogCategoryIds = blogCategoryIds

	file, fileName := handleSingleFile(ctx)
	if err := handler.blogUc.AdjustBlog(ctx, ctx.Params("id"), request, file, fileName); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
	} else {
		messages.SendSuccessResponse(ctx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-BLOG-0002",
		})
	}

	return nil
}

func (handler *BlogHttp) PublishBlog(ctx *fiber.Ctx) error {
	request := &requests.UpsertBlogRequest{
		Title:   ctx.FormValue("title"),
		Content: ctx.FormValue("content"),
	}

	blogCategoryIds := handleBlogCategories(ctx)
	request.BlogCategoryIds = blogCategoryIds

	if err := ctx.BodyParser(&request); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	if formErrors, err := validators.ValidateStruct(ctx, request); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error:      err,
			FormErrors: formErrors,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	file, fileName := handleSingleFile(ctx)
	if err := handler.blogUc.PublishBlog(ctx, ctx.Params("id"), request, file, fileName); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
	} else {
		messages.SendSuccessResponse(ctx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-BLOG-0004",
		})
	}

	return nil
}

func (handler *BlogHttp) UpdateBlog(ctx *fiber.Ctx) error {
	request := &requests.UpsertBlogRequest{
		Title:   ctx.FormValue("title"),
		Content: ctx.FormValue("content"),
	}

	blogCategoryIds := handleBlogCategories(ctx)
	request.BlogCategoryIds = blogCategoryIds

	if err := ctx.BodyParser(&request); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	if formErrors, err := validators.ValidateStruct(ctx, request); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error:      err,
			FormErrors: formErrors,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	file, fileName := handleSingleFile(ctx)
	if err := handler.blogUc.UpdateBlog(ctx, ctx.Params("id"), request, file, fileName); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
	} else {
		messages.SendSuccessResponse(ctx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-BLOG-0002",
		})
	}

	return nil
}

func (handler *BlogHttp) DeleteBlog(ctx *fiber.Ctx) error {
	if err := handler.blogUc.DeleteBlog(ctx, ctx.Params("id")); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
	} else {
		messages.SendSuccessResponse(ctx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-BLOG-0003",
		})
	}

	return nil
}

func (handler *BlogHttp) UpdateSlug(ctx *fiber.Ctx) error {
	request := &requests.UpdateSlugRequest{}
	if err := ctx.BodyParser(request); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
		return nil
	}

	if formErrors, err := validators.ValidateStruct(ctx, request); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error:      err,
			FormErrors: formErrors,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	if err := handler.blogUc.UpdateSlug(ctx, ctx.Params("id"), request); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
	} else {
		messages.SendSuccessResponse(ctx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-BLOG-0005",
		})
	}

	return nil
}

func (handler *BlogHttp) UpdateToDraft(ctx *fiber.Ctx) error {
	if err := handler.blogUc.UpdateBlogToDraft(ctx, ctx.Params("id")); err != nil {
		messages.SendErrorResponse(ctx, responses.ErrorResponse{
			Error: err,
		})
	} else {
		messages.SendSuccessResponse(ctx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-BLOG-0006",
		})
	}

	return nil
}

func handleBlogCategories(ctx *fiber.Ctx) []int {
	blogCategoryIds := ctx.FormValue("blogCategoryIds")
	idStrings := strings.Split(blogCategoryIds, ",")
	var convertedCategoryIds []int
	for _, id := range idStrings {
		convertedId, err := strconv.Atoi(id)
		if err != nil {
			messages.SendErrorResponse(ctx, responses.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Error:      err,
			})
			return nil
		}
		convertedCategoryIds = append(convertedCategoryIds, convertedId)
	}

	return convertedCategoryIds
}

func handleSingleFile(ctx *fiber.Ctx) ([]byte, string) {
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
