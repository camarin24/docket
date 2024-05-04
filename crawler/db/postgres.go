package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(cfg DbConfig) *gorm.DB {
	//TODO: Properly set sslmode
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic("fail to connect database")
	}
	return db
}
