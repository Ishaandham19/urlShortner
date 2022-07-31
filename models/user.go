package models

import (
	"time"
)

//User struct declaration
type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
	UserName  string    `gorm:"uniqueIndex"`
	Email     string    `gorm:"type:varchar(100)"`
	Password  string    `json:"Password"`
}
