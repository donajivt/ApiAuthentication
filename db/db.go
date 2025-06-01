package db

import (
	"github.com/donajivt/go-auth-service/config"
	"github.com/donajivt/go-auth-service/models"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	var err error
	DB, err = gorm.Open(sqlserver.Open(config.Cfg.DSN), &gorm.Config{})
	if err != nil {
		return err
	}

	return DB.AutoMigrate(&models.User{})
}
