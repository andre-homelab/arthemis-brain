package models

import (
	"gorm.io/gorm"
)

type Activity struct {
	gorm.Model
	ProjectID     uint   `gorm:"uniqueIndex"`
	Name          string `gorm:"type:varchar(100)"`
	Description   string `gorm:"type:text"`
	Justification string `gorm:"type:text"`
}
