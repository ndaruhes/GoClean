package factories

import (
	"go-clean/src/app/infrastructures"
	"go-clean/src/domains/blogs/entities"
	masterEntities "go-clean/src/domains/master/entities"
	"math/rand"
)

func BlogCategoryFactory() []entities.BlogCategory {
	var blogCategories []entities.BlogCategory
	blogs := getAllBlogs()
	categoryBlogs := getCategoryBlogs()

	for i := 0; i < 10; i++ {
		randomBlogIdx := rand.Intn(len(blogs))
		randomBlog := blogs[randomBlogIdx]

		randomCategoryBlogIdx := rand.Intn(len(categoryBlogs))
		randomCategoryBlog := categoryBlogs[randomCategoryBlogIdx]

		blogCategories = append(blogCategories, entities.BlogCategory{
			BlogID:         randomBlog.ID,
			CategoryBlogID: randomCategoryBlog.ID,
		})
	}
	return blogCategories
}

func getAllBlogs() []entities.Blog {
	var blogs []entities.Blog
	db := infrastructures.ConnectDatabase()
	db.Model(&entities.Blog{}).Select("*").Scan(&blogs)
	return blogs
}

func getCategoryBlogs() []masterEntities.CategoryBlog {
	var categoryBlogs []masterEntities.CategoryBlog
	db := infrastructures.ConnectDatabase()
	db.Model(&masterEntities.CategoryBlog{}).Select("*").Scan(&categoryBlogs)
	return categoryBlogs
}
