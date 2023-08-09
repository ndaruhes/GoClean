package entities

import (
	user "go-clean/src/domains/users/entities"
	"go-clean/src/shared/entities"
	"strings"
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type Blog struct {
	ID          string     `json:"id" gorm:"primaryKey;type:varchar(255);not null"`
	Title       *string    `json:"title" gorm:"type:varchar(255);null"`
	Slug        string     `json:"slug" gorm:"index:slug,unique,type:varchar(255);not null"`
	Cover       *string    `json:"cover" gorm:"type:varchar(255);null"`
	Content     *string    `json:"content" gorm:"type:text;null"`
	Status      string     `json:"status" gorm:"type:enum('Draft', 'Published');default:'Draft';not null"`
	PublishedAt *time.Time `json:"published_at"`
	UserID      string     `json:"user_id" gorm:"type:varchar(255);not null"`
	User        user.User
	entities.Timestamp
}

func (blog *Blog) BeforeCreate(db *gorm.DB) error {
	if blog.ID == "" {
		blog.ID = strings.ToUpper(xid.New().String())
	}

	return nil
}
