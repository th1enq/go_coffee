package db

import (
	"fmt"

	"github.com/th1enq/go_coffee/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

func LoadDB(cfg *config.Config, serviceType config.ServiceType) (*DB, error) {
	var dbName string
	switch serviceType {
	case config.UserService:
		dbName = cfg.DB.UserDB
	case config.CharacterService:
		dbName = cfg.DB.CharDB
	default:
		return nil, fmt.Errorf("invalid service type: %s", serviceType)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.User,
		cfg.DB.Password,
		dbName,
		cfg.DB.Port)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	return &DB{gormDB}, nil
}
