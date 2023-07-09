package entities

import "go-clean/shared/entities"

type CategoryBlog struct {
	ID               int `json:"id" gorm:"primaryKey"`
	OrderID          int `json:"order_id"`
	CategoryBlogLang []CategoryBlogLang
	entities.Timestamp
}

type CategoryBlogLang struct {
	CategoryBlogID int    `json:"category_blog_id" gorm:"primaryKey"`
	Lang           string `json:"lang" gorm:"primaryKey;type=char(6)"`
	Name           string `json:"name" gorm:"type:varchar(255)"`
	CategoryBlog   CategoryBlog
	entities.Timestamp
}
