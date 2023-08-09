package factories

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/rs/xid"
	"go-clean/src/app/infrastructures"
	"go-clean/src/domains/blogs/entities"
	userEntities "go-clean/src/domains/users/entities"
	"go-clean/src/shared/utils"
	"math/rand"
	"strings"
)

func BlogFactory() []entities.Blog {
	var blogs []entities.Blog
	users := getAllUser()

	for i := 0; i < 10; i++ {
		sentence := faker.Sentence()
		cover := fmt.Sprintf("https://picsum.photos/640/300?random=%d", i)
		content := faker.Paragraph()
		randomUserIdx := rand.Intn(len(users))
		randomUser := users[randomUserIdx]

		blogs = append(blogs, entities.Blog{
			ID:      strings.ToUpper(xid.New().String()),
			Title:   &sentence,
			Slug:    utils.GenerateSlug(sentence),
			Cover:   &cover,
			Content: &content,
			Status:  utils.GetRandomString([]string{"Draft", "Published"}),
			UserID:  randomUser.ID,
		})
	}

	return blogs
}

func getAllUser() []userEntities.User {
	var users []userEntities.User
	db := infrastructures.ConnectDatabase()
	db.Model(&userEntities.User{}).Select("*").Scan(&users)
	return users
}
