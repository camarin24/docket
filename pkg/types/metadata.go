package types

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Metadata struct {
	gorm.Model
	ID         int `gorm:"primaryKey"`
	Name       string
	StorageKey string
	Data       datatypes.JSON
}
