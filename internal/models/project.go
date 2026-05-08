package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	name          string    `gorm:"type:varchar(150);not null"`
	lifetimeStart time.Time `gorm:"type:date"`
	lifetimeEnd   time.Time `gorm:"type:date"`
	justification string    `gorm:"type:text"`
}
