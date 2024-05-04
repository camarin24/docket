package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	ID                uuid.UUID `gorm:"primaryKey,type:uuid"`
	Name              string    `json:"name"`
	StorageKey        string    `json:"storageKey"`
	OriginalPath      string    `json:"originalPath"`
	Size              int64     `json:"size"`
	Content           string    `gorm:"type:text" json:"content"`
	MetaData          string    `json:"metaData"`
	MetaDataExtracted bool      `json:"metaDataExtracted" gorm:"default:false"`
}
