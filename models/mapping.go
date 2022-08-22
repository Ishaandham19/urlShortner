package models

import "time"

//Token struct declaration
type Mapping struct {
	ID             uint
	UserName       string `gorm:"index:idx_member"`
	Alias          string `gorm:"index:idx_member"`
	Url            string
	ExpirationDate time.Time
}
