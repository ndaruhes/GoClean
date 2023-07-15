package requests

type UpsertBlogRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	//BlogCategoryIds []int  `json:"blogCategoryIds" form:"blogCategoryIds[]" validate:"required,exists=category_blogs;id;deleted_at;NULL"`
	BlogCategoryIds []int `json:"blogCategoryIds" form:"blogCategoryIds[]" validate:"required"`
}
