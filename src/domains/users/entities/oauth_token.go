package entities

import (
	"go-clean/shared/entities"
	"strings"
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type OAuthToken struct {
	ID        string    `json:"id" gorm:"primaryKey,type:varchar(255)"`
	User      User      `json:"user"`
	UserID    string    `json:"user_id" gorm:"type:varchar(255)"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked" gorm:"default:0"`
	entities.Timestamp
}

func (token *OAuthToken) BeforeCreate(db *gorm.DB) error {
	token.ID = strings.ToUpper(xid.New().String())
	return nil
}
