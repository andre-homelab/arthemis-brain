package models

import "gorm.io/gorm"

type Proponent struct {
	gorm.Model
	// ID    string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	name  string `gorm:"not null"`
	email string `gorm:"not null"`
}
