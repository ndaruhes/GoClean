package migration

import (
	blogs "go-clean/src/domains/blogs/entities"
	master "go-clean/src/domains/master/entities"
	users "go-clean/src/domains/users/entities"
	"go-clean/src/setup/database"
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
