package responses

import "time"

type PublicBlogListsResponse struct {
	Title       string    `json:"title"`
	Cover       string    `json:"cover"`
	Content     string    `json:"content"`
	Author      string    `json:"author"`
	PublishedAt time.Time `json:"publishedAt"`
}

type BlogDetailResponse struct {
	Title       string    `json:"title"`
	Cover       string    `json:"cover"`
	Content     string    `json:"content"`
	Author      string    `json:"author"`
	PublishedAt time.Time `json:"publishedAt"`
}
