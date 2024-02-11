package db

import (
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(config *viper.Viper) *gorm.DB {

	// Get config
	//host := config.GetString("DB_HOST")
	//port := config.GetString("DB_PORT")
	//user := config.GetString("DB_USER")
	//password := config.GetString("DB_PASSWORD")
	//name := config.GetString("DB_NAME")
	//sslmode := config.GetString("DB_SSLMODE")
	//
	//// Build DSN
	//dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=America/Guayaquil", host, user, password, name, port, sslmode)
	dsn := config.GetString("DB_DSN")
	configDB := &gorm.Config{}
	if config.GetBool("DB_DEBUG") {
		configDB = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
	} else {
		configDB = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}
	}

	db, err := gorm.Open(postgres.Open(dsn), configDB)
	if err != nil {
		panic(err)
	}

	DB = db

	return db
}
