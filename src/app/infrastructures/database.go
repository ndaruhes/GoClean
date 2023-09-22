package infrastructures

import (
	"fmt"
	"go-clean/src/app/config"
	"log"
	"time"

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

		var (
			err        error
			maxRetries = 5
		)

		for i := 1; i <= maxRetries; i++ {
			database, err := gorm.Open(postgres.New(
				postgres.Config{
					DSN: dsn,
				},
			), &gorm.Config{
				Logger:               logger.Default.LogMode(logger.Info),
				FullSaveAssociations: true,
			})

			if err != nil {
				log.Printf("(Attempt %d/%d)\n", i, maxRetries)
				log.Printf("Failed to connect to the database. Retrying in 5 seconds...")
				time.Sleep(5 * time.Second)
			} else {
				db = database
				break
			}
		}

		if err != nil {
			panic(err)
		}
	}

	return db
}
