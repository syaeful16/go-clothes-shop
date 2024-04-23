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
	// mysql://nknk7peybjlh5p0t:f797ey0ny54c4whp@lyl3nln24eqcxxot.cbetxkdyhwsb.us-east-1.rds.amazonaws.com:3306/hmcmb1onm2ixyg2r
	db, err := gorm.Open(mysql.Open("nknk7peybjlh5p0t:f797ey0ny54c4whp@tcp(lyl3nln24eqcxxot.cbetxkdyhwsb.us-east-1.rds.amazonaws.com:3306)/hmcmb1onm2ixyg2r"))
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&User{}, &Product{}, &DetailProduct{}); err != nil {
		panic(err)
	}

	DB = db

}
