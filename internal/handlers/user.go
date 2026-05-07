package handlers

import (
	"gorm.io/gorm"
)

func UserHandler(db *gorm.DB) *DBStore {
	return &DBStore{db}
}

func AddUser() {
	return
}

func DeleteUser() {
	return
}

func EditUser() {
	return
}

func GetUser() {
	return
}
