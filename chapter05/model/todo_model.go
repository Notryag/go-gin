package model

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	ID      int    `gorm:"primaryKey" json:"id"`
	Title   string `gorm:"type:varchar(20);not null" json:"title"`
	Status  string `gorm:"type:varchar(20);" json:"status"`
	Content string `gorm:"type:varchar(20);" json:"content"`
	UserID  int    `json:"userId"`
}
