package types

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Metadata struct {
	gorm.Model
	ID         int `gorm:"primaryKey"`
	DocumentId int
	Data       datatypes.JSON
}
