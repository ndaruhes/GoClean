package http

import (
	"go-clean/src/app/infrastructures"
	"go-clean/src/domains/blogs/interfaces"
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

	"github.com/gofiber/fiber/v2"
)

type BlogHttp struct {
	blogUc interfaces.BlogUseCase
}

func NewBlogHttp(route *fiber.App) *BlogHttp {
	db := infrastructures.ConnectDatabase()
	elasticsearch := infrastructures.CreateElasticsearchClient()
	blogRepo := blogRepository.NewBlogRepository(db, elasticsearch)
	blogUc := blogUseCase.NewBlogUseCase(blogRepo, db)

	handler := &BlogHttp{blogUc: blogUc}
	setRoutes(route, handler)

	return handler
}

func setRoutes(route *fiber.App, handler *BlogHttp) {
	blog := route.Group("blog")
	{
		blog.Get("", handler.GetPublicBlogList)
		blog.Get("/search", handler.SearchBlog)
		blog.Get("/:id", handler.GetBlogDetail)
		blog.Use(middlewares.AuthMiddleware())
		blog.Post("", handler.CreateBlog)
		blog.Put("/drafted/:id/update", handler.UpdateDraftedBlog)
		blog.Put("/:id/publish", handler.PublishBlog)
		blog.Put("/published/:id/update", handler.UpdatePublishedBlog)
		blog.Delete("/:id", handler.DeleteBlog)
		blog.Put("/:id/slug", handler.UpdateSlug)
		blog.Put("/:id/draft", handler.UpdateToDraft)
	}
}

func (handler *BlogHttp) GetPublicBlogList(fiberCtx *fiber.Ctx) error {
	ctx := utils.GetContext(fiberCtx)

	request := &requests.BlogListRequest{}
	if err := fiberCtx.QueryParser(request); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	if formErrors, err := validators.ValidateStruct(ctx, request); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error:      err,
			FormErrors: formErrors,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	if data, totalData, err := handler.blogUc.GetPublicBlogList(ctx, request); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error: err,
		})
	} else {
		messages.SendSuccessResponse(fiberCtx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-BLOG-0007",
			Data:        data,
			TotalData:   &totalData,
		})
	}

	return nil
}

func (handler *BlogHttp) SearchBlog(fiberCtx *fiber.Ctx) error {
	ctx := utils.GetContext(fiberCtx)
	request := &requests.SearchBlogRequest{}
	if err := fiberCtx.QueryParser(request); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	if formErrors, err := validators.ValidateStruct(ctx, request); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error:      err,
			FormErrors: formErrors,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	if data, totalData, err := handler.blogUc.SearchBlog(ctx, request); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error: err,
		})
	} else {
		messages.SendSuccessResponse(fiberCtx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-BLOG-0008",
			Data:        data,
			TotalData:   &totalData,
		})
	}

	return nil
}

func (handler *BlogHttp) GetBlogDetail(fiberCtx *fiber.Ctx) error {
	ctx := utils.GetContext(fiberCtx)
	if data, err := handler.blogUc.GetBlogDetail(ctx, fiberCtx.Params("id")); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error: err,
		})
	} else {
		messages.SendSuccessResponse(fiberCtx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-BLOG-0007",
			Data:        data,
		})
	}

	return nil
}

func (handler *BlogHttp) CreateBlog(fiberCtx *fiber.Ctx) error {
	ctx := utils.GetContext(fiberCtx)
	request := &requests.UpsertBlogRequest{
		Title:   fiberCtx.FormValue("title"),
		Content: fiberCtx.FormValue("content"),
	}

	blogCategoryIds := handleBlogCategories(fiberCtx)
	request.BlogCategoryIds = blogCategoryIds

	file, fileName, err := handleSingleFile(fiberCtx)
	if err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	if err := handler.blogUc.CreateBlog(ctx, request, file, fileName); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error: err,
		})
	} else {
		messages.SendSuccessResponse(fiberCtx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-BLOG-0001",
		})
	}

	return nil
}

func (handler *BlogHttp) UpdateDraftedBlog(fiberCtx *fiber.Ctx) error {
	ctx := utils.GetContext(fiberCtx)
	request := &requests.UpsertBlogRequest{
		Title:   fiberCtx.FormValue("title"),
		Content: fiberCtx.FormValue("content"),
	}

	blogCategoryIds := handleBlogCategories(fiberCtx)
	request.BlogCategoryIds = blogCategoryIds

	file, fileName, err := handleSingleFile(fiberCtx)
	if err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	if err := handler.blogUc.UpdateDraftedBlog(ctx, fiberCtx.Params("id"), request, file, fileName); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error: err,
		})
	} else {
		messages.SendSuccessResponse(fiberCtx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-BLOG-0002",
		})
	}

	return nil
}

func (handler *BlogHttp) PublishBlog(fiberCtx *fiber.Ctx) error {
	ctx := utils.GetContext(fiberCtx)
	request := &requests.UpsertBlogRequest{
		Title:   fiberCtx.FormValue("title"),
		Content: fiberCtx.FormValue("content"),
	}

	blogCategoryIds := handleBlogCategories(fiberCtx)
	request.BlogCategoryIds = blogCategoryIds

	if err := fiberCtx.BodyParser(request); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	if formErrors, err := validators.ValidateStruct(ctx, request); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error:      err,
			FormErrors: formErrors,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	file, fileName, err := handleSingleFile(fiberCtx)
	if err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	if err := handler.blogUc.PublishBlog(ctx, fiberCtx.Params("id"), request, file, fileName); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error: err,
		})
	} else {
		messages.SendSuccessResponse(fiberCtx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-BLOG-0004",
		})
	}

	return nil
}

func (handler *BlogHttp) UpdatePublishedBlog(fiberCtx *fiber.Ctx) error {
	ctx := utils.GetContext(fiberCtx)
	request := &requests.UpsertBlogRequest{
		Title:   fiberCtx.FormValue("title"),
		Content: fiberCtx.FormValue("content"),
	}

	blogCategoryIds := handleBlogCategories(fiberCtx)
	request.BlogCategoryIds = blogCategoryIds

	if err := fiberCtx.BodyParser(request); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	if formErrors, err := validators.ValidateStruct(ctx, request); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error:      err,
			FormErrors: formErrors,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	file, fileName, err := handleSingleFile(fiberCtx)
	if err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	if err := handler.blogUc.UpdatePublishedBlog(ctx, fiberCtx.Params("id"), request, file, fileName); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error: err,
		})
	} else {
		messages.SendSuccessResponse(fiberCtx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-BLOG-0002",
		})
	}

	return nil
}

func (handler *BlogHttp) DeleteBlog(fiberCtx *fiber.Ctx) error {
	ctx := utils.GetContext(fiberCtx)
	if err := handler.blogUc.DeleteBlog(ctx, fiberCtx.Params("id")); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error: err,
		})
	} else {
		messages.SendSuccessResponse(fiberCtx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-BLOG-0003",
		})
	}

	return nil
}

func (handler *BlogHttp) UpdateSlug(fiberCtx *fiber.Ctx) error {
	ctx := utils.GetContext(fiberCtx)
	request := &requests.UpdateSlugRequest{}
	if err := fiberCtx.BodyParser(request); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error: err,
		})
		return nil
	}

	if formErrors, err := validators.ValidateStruct(ctx, request); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error:      err,
			FormErrors: formErrors,
			StatusCode: http.StatusBadRequest,
		})
		return nil
	}

	if err := handler.blogUc.UpdateSlug(ctx, fiberCtx.Params("id"), request); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error: err,
		})
	} else {
		messages.SendSuccessResponse(fiberCtx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-BLOG-0005",
		})
	}

	return nil
}

func (handler *BlogHttp) UpdateToDraft(fiberCtx *fiber.Ctx) error {
	ctx := utils.GetContext(fiberCtx)
	if err := handler.blogUc.UpdateBlogToDraft(ctx, fiberCtx.Params("id")); err != nil {
		messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
			Error: err,
		})
	} else {
		messages.SendSuccessResponse(fiberCtx, responses.SuccessResponse{
			SuccessCode: "SUCCESS-BLOG-0006",
		})
	}

	return nil
}

func handleBlogCategories(fiberCtx *fiber.Ctx) []int {
	blogCategoryIds := fiberCtx.FormValue("blogCategoryIds")
	idStrings := strings.Split(blogCategoryIds, ",")
	var convertedCategoryIds []int
	for _, id := range idStrings {
		convertedId, err := strconv.Atoi(id)
		if err != nil {
			messages.SendErrorResponse(fiberCtx, responses.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Error:      err,
			})
			return nil
		}
		convertedCategoryIds = append(convertedCategoryIds, convertedId)
	}

	return convertedCategoryIds
}

func handleSingleFile(fiberCtx *fiber.Ctx) ([]byte, string, error) {
	var (
		file     []byte
		fileName string
	)
	header, err := fiberCtx.FormFile("cover")
	if header != nil {
		if messages.HasError(err) {
			return nil, "", err
		}

		err = validators.ValidateImage(header)
		if messages.HasError(err) {
			return nil, "", err
		}

		file, err = utils.MultipartFileHeaderToByte(header)
		if messages.HasError(err) {
			return nil, "", err
		}

		fileName = utils.GenerateFileName(header)
	}

	return file, fileName, nil
}
