package config

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(config *Config, log *zap.Logger) *gorm.DB {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.Database.Username, config.Database.Password, config.Database.Address, config.Database.Port, config.Database.Name)

	db, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Error InitDatabase sql open connection fatal error: %v", zap.Error(err))
	}
	db.Logger.LogMode(logger.Info)
	if err = db.Error; err != nil {
		log.Fatal("Error InitDatabase fatal error: ", zap.Error(err))
	}
	log.Info("Connection Success")
	return db
}
