package usecases

import (
	"go-clean/domains/blogs"
	"go-clean/domains/blogs/constants"
	"go-clean/domains/blogs/entities"
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

func (uc *BlogUseCase) GetPublicBlogList(ctx *gin.Context) (*responses.PublicBlogListsResponse, error) {
	//TODO implement me
	panic("implement me")
}

// BLOG USECASE
func (uc *BlogUseCase) CreateBlog(ctx *gin.Context, request *requests.UpsertBlogRequest, file []byte, fileName string) error {
	user := ctx.Value("member").(*responses.TokenDecoded)

	title := "Untitled Blog"
	if request.Title != "" {
		title = request.Title
	}
	blogPayload := &entities.Blog{
		Title:   &title,
		Cover:   &fileName,
		Slug:    utils.GenerateSlug(title),
		Content: &request.Content,
		UserID:  user.ID,
	}

	if _, err := validators.ValidateStruct(ctx, blogPayload); err != nil {
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusBadRequest,
		}
	}

	tx := uc.blogRepo.BeginTransaction(ctx)
	ctx.Set("tx", tx)
	newBlog, err := uc.blogRepo.CreateBlog(ctx, blogPayload)
	if err != nil {
		tx.Rollback()
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	var blogCategoryPayload []entities.BlogCategory
	for _, id := range request.BlogCategoryIds {
		blogCategoryPayload = append(blogCategoryPayload, entities.BlogCategory{
			BlogID:         newBlog.ID,
			CategoryBlogID: id,
		})
	}

	if err := uc.blogRepo.CreateBlogCategory(ctx, blogCategoryPayload); err != nil {
		tx.Rollback()
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	uc.blogRepo.Commit(ctx)

	if file != nil && fileName != "" {
		compressed, err := utils.CompressFile(file, 70)
		if err != nil {
			return err
		}

		imgDir := constants.IMGPATH + user.ID
		err = utils.UploadSingleFile(compressed, imgDir, fileName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *BlogUseCase) AdjustBlog(ctx *gin.Context, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error {
	user := ctx.Value("member").(*responses.TokenDecoded)

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

	if file != nil && fileName != "" {
		payload.Cover = &fileName
	}

	err = validateStatusAndStruct(ctx, blog, constants.DRAFT, payload)
	if err != nil {
		return err
	}

	tx := uc.blogRepo.BeginTransaction(ctx)
	ctx.Set("tx", tx)

	if err := uc.blogRepo.UpdateBlog(ctx, blogID, constants.DRAFT, payload); err != nil {
		tx.Rollback()
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	var blogCategoryPayload []entities.BlogCategory
	for _, id := range request.BlogCategoryIds {
		blogCategoryPayload = append(blogCategoryPayload, entities.BlogCategory{
			BlogID:         blog.ID,
			CategoryBlogID: id,
		})
	}

	if err := uc.blogRepo.UpdateBlogCategory(ctx, blog.ID, blogCategoryPayload); err != nil {
		tx.Rollback()
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	uc.blogRepo.Commit(ctx)

	if file != nil && fileName != "" {
		compressed, err := utils.CompressFile(file, 70)
		if err != nil {
			return err
		}

		imgDir := constants.IMGPATH + user.ID
		err = utils.UploadSingleFile(compressed, imgDir, fileName)
		if err != nil {
			return err
		}

		err = utils.DeleteSingleFile(constants.IMGPATH+user.ID, *blog.Cover)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *BlogUseCase) PublishBlog(ctx *gin.Context, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error {
	user := ctx.Value("member").(*responses.TokenDecoded)

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

	if file != nil && fileName != "" {
		payload.Cover = &fileName
	} else {
		return &messages.ErrorWrapper{
			Context:    ctx,
			ErrorCode:  "ERROR-400004",
			StatusCode: http.StatusBadRequest,
		}
	}

	err = validateStatusAndStruct(ctx, blog, constants.DRAFT, payload)
	if err != nil {
		return err
	}

	tx := uc.blogRepo.BeginTransaction(ctx)
	ctx.Set("tx", tx)

	if err := uc.blogRepo.UpdateBlog(ctx, blogID, constants.DRAFT, payload); err != nil {
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	var blogCategoryPayload []entities.BlogCategory
	for _, id := range request.BlogCategoryIds {
		blogCategoryPayload = append(blogCategoryPayload, entities.BlogCategory{
			BlogID:         blog.ID,
			CategoryBlogID: id,
		})
	}

	if err := uc.blogRepo.UpdateBlogCategory(ctx, blog.ID, blogCategoryPayload); err != nil {
		tx.Rollback()
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	uc.blogRepo.Commit(ctx)

	if file != nil && fileName != "" {
		compressed, err := utils.CompressFile(file, 70)
		if err != nil {
			return err
		}

		imgDir := constants.IMGPATH + user.ID
		err = utils.UploadSingleFile(compressed, imgDir, fileName)
		if err != nil {
			return err
		}

		err = utils.DeleteSingleFile(constants.IMGPATH+user.ID, *blog.Cover)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *BlogUseCase) UpdateBlog(ctx *gin.Context, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error {
	user := ctx.Value("member").(*responses.TokenDecoded)

	blog, err := uc.blogRepo.FindBlogById(ctx, blogID)
	if err != nil {
		return err
	}

	payload := &entities.Blog{
		Title:   &request.Title,
		Content: &request.Content,
	}

	if file != nil && fileName != "" {
		payload.Cover = &fileName
	}

	err = validateStatusAndStruct(ctx, blog, constants.PUBLISHED, payload)
	if err != nil {
		return err
	}

	tx := uc.blogRepo.BeginTransaction(ctx)
	ctx.Set("tx", tx)

	if err := uc.blogRepo.UpdateBlog(ctx, blogID, constants.PUBLISHED, payload); err != nil {
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	var blogCategoryPayload []entities.BlogCategory
	for _, id := range request.BlogCategoryIds {
		blogCategoryPayload = append(blogCategoryPayload, entities.BlogCategory{
			BlogID:         blog.ID,
			CategoryBlogID: id,
		})
	}

	if err := uc.blogRepo.UpdateBlogCategory(ctx, blog.ID, blogCategoryPayload); err != nil {
		tx.Rollback()
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	uc.blogRepo.Commit(ctx)

	if file != nil && fileName != "" {
		compressed, err := utils.CompressFile(file, 70)
		if err != nil {
			return err
		}

		imgDir := constants.IMGPATH + user.ID
		err = utils.UploadSingleFile(compressed, imgDir, fileName)
		if err != nil {
			return err
		}

		err = utils.DeleteSingleFile(constants.IMGPATH+user.ID, *blog.Cover)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *BlogUseCase) UpdateSlug(ctx *gin.Context, blogID string, request *requests.UpsertBlogRequest) error {
	blog, err := uc.blogRepo.FindBlogById(ctx, blogID)
	if err != nil {
		return err
	}

	if blog.Status != constants.PUBLISHED {
		return &messages.ErrorWrapper{
			Context:    ctx,
			ErrorCode:  "ERROR-403001",
			StatusCode: http.StatusForbidden,
		}
	}

	payload := &entities.Blog{
		Slug: utils.GenerateSlug(request.Title),
	}

	if err := uc.blogRepo.UpdateBlog(ctx, blogID, constants.PUBLISHED, payload); err != nil {
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

	if blog.Status != constants.PUBLISHED {
		return &messages.ErrorWrapper{
			Context:    ctx,
			ErrorCode:  "ERROR-403001",
			StatusCode: http.StatusForbidden,
		}
	}

	payload := &entities.Blog{
		Status: constants.DRAFT,
	}

	if err := uc.blogRepo.UpdateBlog(ctx, blogID, constants.PUBLISHED, payload); err != nil {
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

	currImgDir := constants.IMGPATH + blog.UserID
	targetImgDir := constants.IMGPATH + blog.UserID + "/trash"
	if err := utils.MoveSingleFile(currImgDir, targetImgDir, *blog.Cover); err != nil {
		return err
	}

	return uc.blogRepo.DeleteBlog(ctx, blogID)
}

func validateStatusAndStruct(ctx *gin.Context, blog *entities.Blog, wantStatus string, payload *entities.Blog) error {
	if blog.Status != wantStatus {
		return &messages.ErrorWrapper{
			Context:    ctx,
			ErrorCode:  "ERROR-403001",
			StatusCode: http.StatusForbidden,
		}
	}

	if _, err := validators.ValidateStruct(ctx, payload); err != nil {
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusBadRequest,
		}
	}

	return nil
}
