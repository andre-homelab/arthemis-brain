package models

import (
	"time"

	"gorm.io/gorm"
)

type IndicatorObservation struct {
	gorm.Model
	IndicatorID   uint
	ValueBaseline float32   `gorm:"type:decimal"`
	Date          time.Time `gorm:"type:date"`
	Position      JSONB     `gorm:"type:jsonb"`
}
