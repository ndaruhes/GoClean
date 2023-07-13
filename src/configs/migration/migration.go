package migration

import (
	"go-clean/configs/database"
	blogs "go-clean/domains/blogs/entities"
	master "go-clean/domains/master/entities"
	users "go-clean/domains/users/entities"
)

func Migrate() error {
	db := database.ConnectDatabase()
	return db.AutoMigrate(
		// MASTER
		&master.CategoryBlog{}, &master.CategoryBlogLang{},

		// USER
		&users.User{}, &users.OAuthToken{},

		//	BLOG
		&blogs.Blog{}, &blogs.BlogCategory{},
	)
}
