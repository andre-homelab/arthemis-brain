package models

import "gorm.io/gorm"

type Proponent struct {
	gorm.Model
	Name    string  `gorm:"not null"`
	Email   string  `gorm:"not null"`
	Project Project `gorm:"constraint:OnDelete:CASCADE;"`
}
