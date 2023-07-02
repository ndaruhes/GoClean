package entities

import (
	"go-clean/shared/entities"
	"strings"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type User struct {
	ID       string  `json:"id" gorm:"primaryKey;type:varchar(255);not null"`
	Name     string  `json:"name" gorm:"type:varchar(255);not null"`
	Email    string  `json:"email" gorm:"index:email,unique,type:varchar(255);not null"`
	Password string  `json:"password" gorm:"type:varchar(255);not null"`
	GoogleID *string `json:"google_id" gorm:"index:google_id,unique,type:varchar(255);null"`
	Role     string  `json:"role" gorm:"type:enum('Admin', 'Member');default:'Member';not null"`
	entities.Timestamp
}

func (user *User) BeforeCreate(db *gorm.DB) error {
	if user.ID == "" {
		user.ID = strings.ToUpper(xid.New().String())
	}
	return nil
}
