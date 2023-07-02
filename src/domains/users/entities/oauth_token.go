package entities

import (
	"go-clean/shared/entities"
	"strings"
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type OAuthToken struct {
	ID        string    `json:"id" gorm:"primaryKey,type:varchar(255);not null"`
	UserID    string    `json:"user_id" gorm:"type:varchar(255);not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	Revoked   bool      `json:"revoked" gorm:"default:0;not null"`
	User      User      `json:"user"`
	entities.Timestamp
}

func (token *OAuthToken) BeforeCreate(db *gorm.DB) error {
	token.ID = strings.ToUpper(xid.New().String())
	return nil
}
