package models

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ModelDefault struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"update_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"delete_at"`
}

var DB *gorm.DB

func ConnectDB() {
	db, err := gorm.Open(mysql.Open("dev:dev123@tcp(localhost:3306)/clothes_db?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		panic(err)
	}

	DB = db

}
