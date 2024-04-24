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
	driverPg := "postgres://ug104hbg8pa8m:p685c1725fdc9fcabbc116365a89d014591757e403d404d92574e6e48eccd2f23@ceu9lmqblp8t3q.cluster-czrs8kj4isg7.us-east-1.rds.amazonaws.com:5432/dbg6sk4uol36ch"
	db, err := gorm.Open(postgres.Open(driverPg))
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&User{}, &Product{}, &DetailProduct{}); err != nil {
		panic(err)
	}

	DB = db

}
