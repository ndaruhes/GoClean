package requests

type UpsertBlogRequest struct {
	Title   string `json:"title" validate:"required"`
	Cover   string `json:"cover" validate:"required"`
	Content string `json:"content" validate:"required"`
}
