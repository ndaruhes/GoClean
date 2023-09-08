package utils

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetDb(ctx *fiber.Ctx, db *gorm.DB) *gorm.DB {
	tx := ctx.Locals("tx")
	if tx != nil {
		return tx.(*gorm.DB)
	}
	return db
}

func BeginTransaction(ctx *fiber.Ctx, db *gorm.DB) *gorm.DB {
	tx := db.Begin()
	ctx.Locals("tx", tx)
	return tx
}

func Commit(ctx *fiber.Ctx) {
	tx := ctx.Locals("tx")
	if tx != nil {
		db := tx.(*gorm.DB)
		db.Commit()
	}
}
