package factories

import (
	"github.com/bxcodec/faker/v3"
	"go-clean/src/domains/master/entities"
)

func CategoryBlogFactory() []entities.CategoryBlog {
	var categoryBlogs []entities.CategoryBlog
	for i := 1; i <= 10; i++ {
		categoryBlogs = append(categoryBlogs, entities.CategoryBlog{
			ID:      i,
			OrderID: i,
		})
	}
	return categoryBlogs
}

func CategoryBlogLangFactory() []entities.CategoryBlogLang {
	var categoryBlogLang []entities.CategoryBlogLang
	for i := 1; i <= 10; i++ {
		categoryBlogLang = append(categoryBlogLang, entities.CategoryBlogLang{
			CategoryBlogID: i,
			Lang:           "en",
			Name:           faker.Word(),
		})
	}
	return categoryBlogLang
}
