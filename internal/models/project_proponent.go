package models

import "gorm.io/gorm"

type PersonAddress struct {
	gorm.Model
	ProponentID Proponent
	ProjectID   Project
	Role        string `gorm:"type:varchar(50)"`
}
