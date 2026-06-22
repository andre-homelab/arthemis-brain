package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ProponentID       uint
	Name              string    `gorm:"type:varchar(150);not null"`
	LifetimeStart     time.Time `gorm:"type:date"`
	LifetimeEnd       time.Time `gorm:"type:date"`
	Justification     string    `gorm:"type:text"`
	Locations         []Location
	Activities        []Activity
	ProjectProponents []ProjectProponent
	ProjectSdgs       []Sdg `gorm:"many2many:project_sdg"`
}

type ProjectRequest struct {
	ProponentID   uint
	Name          string
	LifetimeStart time.Time
	LifetimeEnd   time.Time
	Justification string
	SdgIDs        []uint
}
