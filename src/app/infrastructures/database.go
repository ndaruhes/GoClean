package infrastructures

import (
	"fmt"
	"go-clean/src/app/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func ConnectDatabase() *gorm.DB {
	var dsn string
	if db == nil {
		conf := config.GetConfig().Database
		if conf.Password == "" {
			dsn = fmt.Sprintf("%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Username, conf.Host, conf.Port, conf.Name)
		} else {
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Username, conf.Password, conf.Host, conf.Port, conf.Name)
		}
		database, err := gorm.Open(mysql.New(
			mysql.Config{
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
