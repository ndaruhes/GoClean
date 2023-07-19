package utils

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDb(ctx *gin.Context, db *gorm.DB) *gorm.DB {
	tx := ctx.Value("tx")
	if tx != nil {
		return tx.(*gorm.DB)
	}
	return db
}

func BeginTransaction(ctx *gin.Context, db *gorm.DB) *gorm.DB {
	tx := db.Begin()
	ctx.Set("tx", tx)
	return tx
}

func Commit(ctx *gin.Context) {
	tx := ctx.Value("tx")
	if tx != nil {
		db := tx.(*gorm.DB)
		db.Commit()
	}
}
