package models

import "time"

//Token struct declaration
type Mapping struct {
	ID             uint
	UserName       string
	Alias          string `gorm:"uniqueIndex"`
	Url            string
	ExpirationDate time.Time
}
