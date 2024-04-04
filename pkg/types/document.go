package types

import (
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	ID           int `gorm:"primaryKey"`
	Name         string
	StorageKey   string
	OriginalPath string
	Size         int64
}
