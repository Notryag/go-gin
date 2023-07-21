package database

import (
	"chapter05/config"
	"chapter05/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	config.Init()

	db, err := gorm.Open(mysql.Open(config.Cfg.Dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
		return nil, err
	}
	err = db.AutoMigrate(&model.User{}, &model.Todo{})
	if err != nil {
		panic("failed to migrate database")
		return nil, err
	}

	return db, nil
}
