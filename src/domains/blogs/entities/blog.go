package entities

import (
	"github.com/rs/xid"
	user "go-clean/domains/users/entities"
	"go-clean/shared/entities"
	"gorm.io/gorm"
	"strings"
)

type Blog struct {
	ID      string `json:"id" gorm:"primaryKey;type:varchar(255);not null"`
	Title   string `json:"title" gorm:"type:varchar(255);not null"`
	Slug    string `json:"slug" gorm:"index:slug,unique,type:varchar(255);not null"`
	Cover   string `json:"cover" gorm:"type:varchar(255);not null"`
	Content string `json:"content" gorm:"type:varchar(255);not null"`
	UserID  string `json:"user_id" gorm:"type:varchar(255);not null"`
	User    user.User
	entities.Timestamp
}

func (blog *Blog) BeforeCreate(db *gorm.DB) error {
	if blog.ID == "" {
		blog.ID = strings.ToUpper(xid.New().String())
	}

	return nil
}
