package models

import "gorm.io/gorm"

type ProjectSdg struct {
	gorm.Model
	ProjectID uint
	SdgID     uint
}
