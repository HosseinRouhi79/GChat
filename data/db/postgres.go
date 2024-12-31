package db

import (
	"chat/config"
	"chat/src/api/pkg/logging"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type single struct{}

var cfg = config.GetConfig()

var zap = logging.NewLogger(cfg)

var singleton *single

var dbClient *gorm.DB

func InitDB(cfg *config.Config) error {
	if singleton == nil {
		var err error
		cnn := fmt.Sprintf("host=%s port=%s password=%s user=%s sslmode=%s dbname=%s timezone=Asia/Tehran",
			cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Password, cfg.Postgres.User, cfg.Postgres.SSLMode, cfg.Postgres.DbName)
		//need to open a connection
		dbClient, err = gorm.Open(postgres.Open(cnn), &gorm.Config{})
		if err != nil {
			return err
		}
		sqlDB, err := dbClient.DB()
		if err != nil {
			return err
		}
		err = sqlDB.Ping()
		if err != nil {
			return err
		}
		sqlDB.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(cfg.Postgres.ConnMaxLifetime)

		singleton = &single{}
		zap.Info(logging.Postgres, logging.Migration, "Postgres connected successfully", nil)
	} else {
		zap.Info(logging.Postgres, logging.Migration, "Postgres is already configured", nil)
	}

	return nil
}

func GetDB() *gorm.DB {
	return dbClient
}

func CloseDB() error {
	cnn, _ := dbClient.DB()
	cnn.Close()
	return nil
}
