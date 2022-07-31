package models

import "gorm.io/gorm"

//User struct declaration
type User struct {
	gorm.Model
	UserName string `gorm:"uniqueIndex"`
	Email    string `gorm:"type:varchar(100)"`
	Password string `json:"Password"`
}
