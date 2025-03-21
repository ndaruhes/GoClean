package operation

import (
	"context"
	"fmt"
	"go-clean/src/models/requests"
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

func PaginateOrder(request requests.PaginationRequest) func(db *gorm.DB) *gorm.DB {
	page := request.Page
	size := request.Size
	sortOrder := request.SortOrder
	sortBy := request.SortBy

	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * size
		db = db.Offset(offset).Limit(size)

		if sortOrder != "" && sortBy != "" {
			db = db.Order(fmt.Sprintf("%s %s", sortOrder, sortBy))
		}

		return db
	}
}
