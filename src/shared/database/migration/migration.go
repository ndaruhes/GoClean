package migration

import (
	"go-clean/src/app/infrastructures"
	blogs "go-clean/src/domains/blogs/entities"
	master "go-clean/src/domains/master/entities"
	users "go-clean/src/domains/users/entities"
)

func Migrate() error {
	db := infrastructures.ConnectDatabase()
	return db.AutoMigrate(
		// MASTER
		&master.CategoryBlog{}, &master.CategoryBlogLang{},

		// USER
		&users.User{}, &users.OAuthToken{},

		//	BLOG
		&blogs.Blog{}, &blogs.BlogCategory{},
	)
}
