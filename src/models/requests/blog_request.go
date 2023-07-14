package requests

type UpsertBlogRequest struct {
	Title           string `json:"title" validate:"required"`
	Content         string `json:"content" validate:"required"`
	BlogCategoryIds []int  `json:"blogCategoryIds" form:"blogCategoryIds[]" validate:"exists=category_blogs;deleted_at;NULL"`
}
