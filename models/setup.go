package models

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ModelDefault struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"update_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"delete_at"`
}

var DB *gorm.DB

func ConnectDB() {
	// dev:dev123@tcp(localhost:3306)/clothes_db?charset=utf8mb4&parseTime=True&loc=Local
	// mysql://nknk7peybjlh5p0t:f797ey0ny54c4whp@lyl3nln24eqcxxot.cbetxkdyhwsb.us-east-1.rds.amazonaws.com:3306/hmcmb1onm2ixyg2r
	// db, err := gorm.Open(mysql.Open("dev:dev123@tcp(localhost:3306)/clothes_db?charset=utf8mb4&parseTime=True&loc=Local"))
	driverPg := "postgres://syaefulloharnas:040818@localhost:5432/clothes_shop_db"
	db, err := gorm.Open(postgres.Open(driverPg))
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&User{}, &Product{}, &DetailProduct{}); err != nil {
		panic(err)
	}

	DB = db

}
