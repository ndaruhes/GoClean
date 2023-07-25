package entities

import (
	master "go-clean/src/domains/master/entities"
)

type BlogCategory struct {
	BlogID         string `json:"blog_id"`
	CategoryBlogID int    `json:"category_blog_id"`
	Blog           Blog
	CategoryBlog   master.CategoryBlog
}
