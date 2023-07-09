package entities

import (
	master "go-clean/domains/master/entities"
)

type BlogCategory struct {
	BlogID         string `json:"blog_id"`
	CategoryBlogID string `json:"category_blog_id"`
	Blog           Blog
	Category       master.CategoryBlog
}
