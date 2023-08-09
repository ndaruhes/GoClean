package seeder

import (
	"go-clean/src/app/infrastructures"
	"go-clean/src/shared/database/factories"
)

type Seeder struct {
	Seeder interface{}
}

func RegisterSeeders() []Seeder {
	return []Seeder{
		{Seeder: factories.UserFactory()},
		{Seeder: factories.BlogFactory()},
	}
}

func DBSeed() error {
	db := infrastructures.ConnectDatabase()
	for _, seeder := range RegisterSeeders() {
		err := db.Debug().Create(seeder.Seeder).Error
		if err != nil {
			return err
		}
	}

	return nil
}
