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
		{Factory: factories.UserFactory()},
		{Factory: factories.BlogFactory()},
	}
}

func DBSeed() error {
	db := infrastructures.ConnectDatabase()
	for _, seeder := range RegisterSeeders() {
		err := db.Debug().Create(seeder.Factory).Error
		if err != nil {
			return err
		}
	}

	return nil
}
