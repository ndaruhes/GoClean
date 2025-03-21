package seeder

import (
	"go-clean/src/app/infrastructures"
	"go-clean/src/shared/database/factories"
)

type Seeder struct {
	Factory interface{}
}

func RegisterSeeders() []Seeder {
	return []Seeder{
		// MASTER
		{Factory: factories.CategoryBlogFactory()},
		{Factory: factories.CategoryBlogLangFactory()},

		// USER
		{Factory: factories.UserFactory()},

		// BLOG
		{Factory: factories.BlogFactory()},
		{Factory: factories.BlogCategoryFactory()},
	}
}

func DBSeed() error {
	db := infrastructures.ConnectDatabase()
	err := db.Debug().Create(factories.CategoryBlogFactory()).Error
	if err != nil {
		return err
	}

	err = db.Debug().Create(factories.CategoryBlogLangFactory()).Error
	if err != nil {
		return err
	}

	err = db.Debug().Create(factories.UserFactory()).Error
	if err != nil {
		return err
	}

	err = db.Debug().Create(factories.BlogFactory()).Error
	if err != nil {
		return err
	}

	err = db.Debug().Create(factories.BlogCategoryFactory()).Error
	if err != nil {
		return err
	}

	return nil
}
