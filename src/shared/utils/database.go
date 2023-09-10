package utils

import (
	"context"
	"gorm.io/gorm"
)

func GetDb(ctx context.Context, db *gorm.DB) *gorm.DB {
	tx := ctx.Value("tx")
	if tx != nil {
		return tx.(*gorm.DB)
	}
	return db
}

func BeginTransaction(ctx context.Context, db *gorm.DB) (*gorm.DB, context.Context) {
	tx := db.Begin()
	ctx = context.WithValue(ctx, "tx", tx)
	return tx, ctx
}

func Commit(ctx context.Context) {
	tx := ctx.Value("tx")
	if tx != nil {
		db := tx.(*gorm.DB)
		db.Commit()
	}
}
