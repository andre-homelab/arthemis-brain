package models

import "gorm.io/gorm"

type Proponent struct {
	gorm.Model
	// ID    string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name  string `gorm:"not null"`
	Email string `gorm:"not null"`
}
