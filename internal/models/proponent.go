package models

import "gorm.io/gorm"

type Proponent struct {
	gorm.Model
	Name              string    `gorm:"not null"`
	Email             string    `gorm:"not null"`
	Projects          []Project `gorm:"constraint:OnDelete:CASCADE;"`
	ProjectProponents []Project `gorm:"many2many:proect_proponent"`
}
