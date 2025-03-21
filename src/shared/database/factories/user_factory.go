package factories

import (
	"github.com/bxcodec/faker/v3"
	"github.com/rs/xid"
	"go-clean/src/domains/users/entities"
	"go-clean/src/shared/helpers"
	"strings"
)

func UserFactory() []entities.User {
	var users []entities.User
	for i := 0; i < 10; i++ {
		users = append(users, entities.User{
			ID:       strings.ToUpper(xid.New().String()),
			Name:     faker.Name(),
			Email:    faker.Email(),
			Password: generatePassword("12345"),
			Role:     "Member",
		})
	}
	return users
}

func generatePassword(password string) string {
	hashPassword, err := helpers.GeneratePassword(password)
	if err != nil {
		return ""
	}

	return hashPassword
}
