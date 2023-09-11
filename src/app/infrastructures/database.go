package infrastructures

import (
	"fmt"
	"go-clean/src/app/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func ConnectDatabase() *gorm.DB {
	var dsn string
	if db == nil {
		conf := config.GetConfig().Database
		if conf.Password == "" {
			dsn = fmt.Sprintf("host=%s user=%s dbname=%s port=%d sslmode=disable", conf.Host, conf.Username, conf.Name, conf.Port)
		} else {
			dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", conf.Host, conf.Username, conf.Password, conf.Name, conf.Port)
		}
		database, err := gorm.Open(postgres.New(
			postgres.Config{
				DSN: dsn,
			},
		), &gorm.Config{
			Logger:               logger.Default.LogMode(logger.Info),
			FullSaveAssociations: true,
		})
		if err != nil {
			panic(err)
		}

		db = database
	}

	return db
}
