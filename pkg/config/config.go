package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Port                  string
	Database              Database
	AppEnvironment        string
	ProductionEnvironment string
}

type Database struct {
	Username string
	Password string
	Address  string
	Port     string
	Name     string
}

var config Config

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("ERROR .env Not found")
	}
	config.Port = os.Getenv("PORT")
	config.ProductionEnvironment = os.Getenv("PRODUCTION_ENV")
	config.Database.Username = os.Getenv("DB_USERNAME")
	config.Database.Password = os.Getenv("DB_PASSWORD")
	config.Database.Address = os.Getenv("DB_ADDRESS")
	config.Database.Port = os.Getenv("DB_PORT")
	config.Database.Name = os.Getenv("DB_NAME")
	config.AppEnvironment = os.Getenv("APP_ENV")
}

func InitDatabase(username, password, address, port, databaseName string) *gorm.DB {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		username, password, address, port, databaseName)
	db, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Error InitDatabase sql open connection fatal error: %v", err)
	}
	db.Logger.LogMode(logger.Info)
	if err = db.Error; err != nil {
		log.Fatalf("Error InitDatabase fatal error: %v", err)
	}
	log.Print("Connection Success")
	return db
}

func GetConfig() *Config {
	return &config
}
