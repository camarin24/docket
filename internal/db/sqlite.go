package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqlite(cfg DbConfig) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(cfg.DbName), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}
