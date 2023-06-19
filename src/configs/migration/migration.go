package migration

import (
	"go-clean/configs/database"
	users "go-clean/domains/users/entities"
)

func Migrate() error {
	db := database.ConnectDatabase()
	return db.AutoMigrate(
		&users.User{}, &users.OAuthToken{},
	)
}
