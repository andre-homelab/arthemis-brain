package models

import (
	"gorm.io/gorm"
)

type Indicator struct {
	gorm.Model
	ActivityID        uint
	LocationID        uint
	Name              string  `gorm:"type:varchar(100)"`
	Unit              string  `gorm:"type:varchar(50)"`
	ValueBaseline     float32 `gorm:"type:decimal"`
	ValueReference    float32 `gorm:"type:decimal"`
	ObservationMethod string  `gorm:"type:text"`
	Justification     string  `gorm:"type:text"`
	Observations      []Observation
}

type IndicatorRequest struct {
	ActivityID        uint
	LocationID        uint
	Name              string
	Unit              string
	ValueBaseline     float32
	ValueReference    float32
	ObservationMethod string
	Justification     string
}
