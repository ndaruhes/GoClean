package usecases

import (
	"context"
	"go-clean/src/domains/blogs/constants"
	"go-clean/src/domains/blogs/entities"
	"go-clean/src/domains/blogs/interfaces"
	"go-clean/src/models/messages"
	"go-clean/src/models/requests"
	"go-clean/src/models/responses"
	"go-clean/src/shared/utils"
	"go-clean/src/shared/validators"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type BlogUseCase struct {
	blogRepo interfaces.BlogRepository
	db       *gorm.DB
}

func NewBlogUseCase(blogRepo interfaces.BlogRepository, db *gorm.DB) *BlogUseCase {
	return &BlogUseCase{
		blogRepo: blogRepo,
		db:       db,
	}
}

func (uc *BlogUseCase) GetPublicBlogList(ctx context.Context) ([]responses.PublicBlogListsResponse, error) {
	data, err := uc.blogRepo.GetPublicBlogList(ctx)
	if err != nil {
		return nil, err
	}
	var blogs []responses.PublicBlogListsResponse

	for _, blog := range data {
		blogResponses := responses.PublicBlogListsResponse{
			Title:       utils.GetStringPointerValue(blog.Title),
			Cover:       utils.GetStringPointerValue(blog.Cover),
			Content:     utils.GetStringPointerValue(blog.Content),
			Author:      blog.User.Name,
			PublishedAt: utils.GetTimePointerValue(blog.PublishedAt),
		}

		blogs = append(blogs, blogResponses)
	}

	return blogs, nil
}

// BLOG USECASE
func (uc *BlogUseCase) CreateBlog(ctx context.Context, request *requests.UpsertBlogRequest, file []byte, fileName string) error {
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

	tx, ctx := utils.BeginTransaction(ctx, uc.db)
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

	utils.Commit(ctx)

	if file != nil && fileName != "" {
		compressed, err := utils.CompressFile(file, 70)
		if err != nil {
			return err
		}

		imgDir := constants.IMG_PATH + user.ID
		err = utils.UploadSingleFile(compressed, imgDir, fileName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *BlogUseCase) AdjustBlog(ctx context.Context, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error {
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

	tx, ctx := utils.BeginTransaction(ctx, uc.db)

	if err := uc.blogRepo.UpdateBlog(ctx, blogID, constants.DRAFT, payload); err != nil {
		tx.Rollback()
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if err := updateBlogCategory(ctx, blog, request, uc, tx); err != nil {
		return err
	}

	utils.Commit(ctx)

	if err = uploadAndDeleteSingleFile(ctx, blog, file, fileName); err != nil {
		return err
	}

	return nil
}

func (uc *BlogUseCase) PublishBlog(ctx context.Context, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error {
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

	tx, ctx := utils.BeginTransaction(ctx, uc.db)

	if err := uc.blogRepo.UpdateBlog(ctx, blogID, constants.DRAFT, payload); err != nil {
		tx.Rollback()
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if err := updateBlogCategory(ctx, blog, request, uc, tx); err != nil {
		return err
	}

	utils.Commit(ctx)

	if err = uploadAndDeleteSingleFile(ctx, blog, file, fileName); err != nil {
		return err
	}

	return nil
}

func (uc *BlogUseCase) UpdateBlog(ctx context.Context, blogID string, request *requests.UpsertBlogRequest, file []byte, fileName string) error {
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

	tx, ctx := utils.BeginTransaction(ctx, uc.db)

	if err := uc.blogRepo.UpdateBlog(ctx, blogID, constants.PUBLISHED, payload); err != nil {
		tx.Rollback()
		return &messages.ErrorWrapper{
			Context:    ctx,
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if err := updateBlogCategory(ctx, blog, request, uc, tx); err != nil {
		return err
	}

	utils.Commit(ctx)

	if err = uploadAndDeleteSingleFile(ctx, blog, file, fileName); err != nil {
		return err
	}

	return nil
}

func (uc *BlogUseCase) UpdateSlug(ctx context.Context, blogID string, request *requests.UpdateSlugRequest) error {
	blog, err := uc.blogRepo.FindBlogById(ctx, blogID)
	if err != nil {
		return err
	}

	payload := &entities.Blog{
		Slug: utils.GenerateSlug(request.Title),
	}

	err = validateStatusAndStruct(ctx, blog, constants.PUBLISHED, payload)
	if err != nil {
		return err
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

func (uc *BlogUseCase) UpdateBlogToDraft(ctx context.Context, blogID string) error {
	blog, err := uc.blogRepo.FindBlogById(ctx, blogID)
	if err != nil {
		return err
	}

	payload := &entities.Blog{
		Status: constants.DRAFT,
	}

	err = validateStatusAndStruct(ctx, blog, constants.PUBLISHED, payload)
	if err != nil {
		return err
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

func (uc *BlogUseCase) DeleteBlog(ctx context.Context, blogID string) error {
	blog, err := uc.blogRepo.FindBlogById(ctx, blogID)
	if err != nil {
		return err
	}

	currImgDir := constants.IMG_PATH + blog.UserID
	targetImgDir := constants.IMG_PATH + blog.UserID + "/trash"
	if err := utils.MoveSingleFile(currImgDir, targetImgDir, *blog.Cover); err != nil {
		return err
	}

	return uc.blogRepo.DeleteBlog(ctx, blogID)
}

func validateStatusAndStruct(ctx context.Context, blog *entities.Blog, wantStatus string, payload *entities.Blog) error {
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

func updateBlogCategory(ctx context.Context, blog *entities.Blog, request *requests.UpsertBlogRequest, uc *BlogUseCase, tx *gorm.DB) error {
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

	return nil
}

func uploadAndDeleteSingleFile(ctx context.Context, blog *entities.Blog, file []byte, fileName string) error {
	user := ctx.Value("member").(*responses.TokenDecoded)

	if file != nil && fileName != "" {
		compressed, err := utils.CompressFile(file, 70)
		if err != nil {
			return err
		}

		imgDir := constants.IMG_PATH + user.ID
		err = utils.UploadSingleFile(compressed, imgDir, fileName)
		if err != nil {
			return err
		}

		err = utils.DeleteSingleFile(constants.IMG_PATH+user.ID, *blog.Cover)
		if err != nil {
			return err
		}
	}

	return nil
}
