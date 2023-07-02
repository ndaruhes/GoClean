package migration

import (
	"go-clean/configs/database"
	blogs "go-clean/domains/blogs/entities"
	users "go-clean/domains/users/entities"
)

func Migrate() error {
	db := database.ConnectDatabase()
	return db.AutoMigrate(
		// USER
		&users.User{}, &users.OAuthToken{},

		//	BLOG
		&blogs.Blog{},
	)
}
