package types

import (
	"time"

	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	ID         int `gorm:"primaryKey"`
	Name       string
	AddedAt    time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	StorageKey string
}


