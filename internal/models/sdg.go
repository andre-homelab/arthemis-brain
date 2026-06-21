package models

import "gorm.io/gorm"

type Sdg struct {
	gorm.Model
	Number  int    `gorm:"type:integer"`
	Name    string `gorm:"type:varchar(150)"`
	IconURL string `gorm:"type:varchar(255)"`
}
