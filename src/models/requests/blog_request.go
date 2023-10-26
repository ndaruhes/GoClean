package requests

type UpsertBlogRequest struct {
	Title           string `json:"title" validate:"required"`
	Content         string `json:"content" validate:"required"`
	BlogCategoryIds []int  `json:"blogCategoryIds" form:"blogCategoryIds[]" validate:"required"`
}

type UpdateSlugRequest struct {
	Title string `json:"title" validate:"required"`
}

type BlogListRequest struct {
	PaginationRequest
}

type SearchBlogRequest struct {
	Keyword string `json:"keyword" form:"keyword"`
	PaginationRequest
}
