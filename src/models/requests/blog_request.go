package requests

type UpsertBlogRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}
