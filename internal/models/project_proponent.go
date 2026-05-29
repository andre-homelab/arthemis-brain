package models

import "gorm.io/gorm"

type ProjectProponent struct {
	gorm.Model
	ProponentID Proponent
	ProjectID   Project
	Role        string `gorm:"type:varchar(50)"`
}
