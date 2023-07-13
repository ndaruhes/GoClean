package entities

import (
	"github.com/rs/xid"
	user "go-clean/domains/users/entities"
	"go-clean/shared/entities"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Blog struct {
	ID          string     `json:"id" gorm:"primaryKey;type:varchar(255);not null"`
	Title       *string    `json:"title" gorm:"type:varchar(255);null"`
	Slug        string     `json:"slug" gorm:"index:slug,unique,type:varchar(255);not null"`
	Cover       *string    `json:"cover" gorm:"type:varchar(255);null"`
	Content     *string    `json:"content" gorm:"type:varchar(255);null"`
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
