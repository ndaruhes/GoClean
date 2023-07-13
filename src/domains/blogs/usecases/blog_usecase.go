package usecases

import (
	"go-clean/domains/blogs"
	"go-clean/domains/blogs/entities"
	"go-clean/domains/blogs/enums"
	"go-clean/models/messages"
	"go-clean/models/requests"
	"go-clean/models/responses"
	"go-clean/shared/utils"
	"go-clean/shared/validators"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BlogUseCase struct {
	blogRepo blogs.BlogRepository
}

func NewBlogUseCase(blogRepo blogs.BlogRepository) *BlogUseCase {
	return &BlogUseCase{
		blogRepo: blogRepo,
	}
}

var imgPath = "public/images/blogs/USER-"

func (uc *BlogUseCase) GetPublicBlogList(ctx *gin.Context) (*responses.PublicBlogListsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (uc *BlogUseCase) CreateBlog(ctx *gin.Context, request *requests.UpsertBlogRequest, file []byte, fileName string) error {
	user := ctx.Value("member").(*responses.TokenDecoded)

	title := "Untitled Blog"
	if request.Title != "" {
		title = request.Title
	}
	payload := &entities.Blog{
		Title:   &title,
		Cover:   &fileName,
		Slug:    utils.GenerateSlug(title),
		Content: &request.Content,
		UserID:  user.ID,
	}

	if _, err := validators.ValidateStruct(ctx, payload); err != nil {
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusBadRequest,
		}
	}

	if file != nil && fileName != "" {
		compressed, err := utils.CompressFile(file, 70)
		if err != nil {
			return err
		}

		imgDir := imgPath + user.ID
		err = utils.UploadSingleFile(compressed, imgDir, fileName)
		if err != nil {
			return err
		}
	}

	if err := uc.blogRepo.CreateBlog(ctx, payload); err != nil {
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (uc *BlogUseCase) AdjustBlog(ctx *gin.Context, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error {
	blog, err := uc.blogRepo.FindBlogById(ctx, blogID)
	if err != nil {
		return err
	}

	if request.Title != "" {
		blog.Title = &request.Title
		blog.Slug = utils.GenerateSlug(request.Title)
	}

	if request.Content != "" {
		blog.Content = &request.Content
	}

	payload := &entities.Blog{
		Title:   blog.Title,
		Slug:    blog.Slug,
		Content: blog.Content,
	}

	payload, err = validateAndUploadSingleFile(ctx, blog, enums.DRAFT, payload, file, fileName)
	if err != nil {
		return err
	}

	if err := uc.blogRepo.UpdateBlog(ctx, blogID, enums.DRAFT, payload); err != nil {
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (uc *BlogUseCase) PublishBlog(ctx *gin.Context, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error {
	blog, err := uc.blogRepo.FindBlogById(ctx, blogID)
	if err != nil {
		return err
	}

	publishTime := time.Now().UTC()
	payload := &entities.Blog{
		Title:       &request.Title,
		Slug:        utils.GenerateSlug(request.Title),
		Content:     &request.Content,
		Status:      "Published",
		PublishedAt: &publishTime,
	}

	payload, err = validateAndUploadSingleFile(ctx, blog, enums.DRAFT, payload, file, fileName)
	if err != nil {
		return err
	}

	if err := uc.blogRepo.UpdateBlog(ctx, blogID, enums.DRAFT, payload); err != nil {
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (uc *BlogUseCase) UpdateBlog(ctx *gin.Context, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error {
	blog, err := uc.blogRepo.FindBlogById(ctx, blogID)
	if err != nil {
		return err
	}

	payload := &entities.Blog{
		Title:   &request.Title,
		Content: &request.Content,
	}

	payload, err = validateAndUploadSingleFile(ctx, blog, enums.PUBLISHED, payload, file, fileName)
	if err != nil {
		return err
	}

	if err := uc.blogRepo.UpdateBlog(ctx, blogID, enums.PUBLISHED, payload); err != nil {
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (uc *BlogUseCase) UpdateSlug(ctx *gin.Context, blogID string, request *requests.UpsertBlogRequest) error {
	blog, err := uc.blogRepo.FindBlogById(ctx, blogID)
	if err != nil {
		return err
	}

	if blog.Status != enums.PUBLISHED {
		return &messages.ErrorWrapper{
			Context:    ctx,
			ErrorCode:  "ERROR-403001",
			StatusCode: http.StatusForbidden,
		}
	}

	payload := &entities.Blog{
		Slug: utils.GenerateSlug(request.Title),
	}

	if err := uc.blogRepo.UpdateBlog(ctx, blogID, enums.PUBLISHED, payload); err != nil {
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (uc *BlogUseCase) UpdateBlogToDraft(ctx *gin.Context, blogID string) error {
	blog, err := uc.blogRepo.FindBlogById(ctx, blogID)
	if err != nil {
		return err
	}

	if blog.Status != enums.PUBLISHED {
		return &messages.ErrorWrapper{
			Context:    ctx,
			ErrorCode:  "ERROR-403001",
			StatusCode: http.StatusForbidden,
		}
	}

	payload := &entities.Blog{
		Status: enums.DRAFT,
	}

	if err := uc.blogRepo.UpdateBlog(ctx, blogID, enums.PUBLISHED, payload); err != nil {
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (uc *BlogUseCase) DeleteBlog(ctx *gin.Context, blogID string) error {
	blog, err := uc.blogRepo.FindBlogById(ctx, blogID)
	if err != nil {
		return err
	}

	currImgDir := imgPath + blog.UserID
	targetImgDir := imgPath + blog.UserID + "/trash"
	if err := utils.MoveSingleFile(currImgDir, targetImgDir, *blog.Cover); err != nil {
		return err
	}

	return uc.blogRepo.DeleteBlog(ctx, blogID)
}

func validateAndUploadSingleFile(ctx *gin.Context, blog *entities.Blog, wantStatus string, payload *entities.Blog, file []byte, fileName string) (*entities.Blog, error) {
	user := ctx.Value("member").(*responses.TokenDecoded)

	if blog.Status != wantStatus {
		return nil, &messages.ErrorWrapper{
			Context:    ctx,
			ErrorCode:  "ERROR-403001",
			StatusCode: http.StatusForbidden,
		}
	}

	if _, err := validators.ValidateStruct(ctx, payload); err != nil {
		return nil, &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusBadRequest,
		}
	}

	if file != nil && fileName != "" {
		payload.Cover = &fileName
		compressed, err := utils.CompressFile(file, 70)
		if err != nil {
			return nil, err
		}

		imgDir := imgPath + user.ID
		err = utils.UploadSingleFile(compressed, imgDir, fileName)
		if err != nil {
			return nil, err
		}

		err = utils.DeleteSingleFile(imgDir, *blog.Cover)
		if err != nil {
			return nil, err
		}
	}

	return payload, nil
}
