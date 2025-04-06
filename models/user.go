package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserID   uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Name     string
	Email    string `gorm:"unique"`
	Password string
}
