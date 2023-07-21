package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       int    `gorm:"primary_key" json:"id"`
	UserName string `gorm:"type:varchar(20);not null" json:"userName"`
	Password string `gorm:"type:varchar(220);not null" json:"password"`
}
