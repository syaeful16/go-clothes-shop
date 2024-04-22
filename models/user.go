package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"type:varchar(255); not null; unique" validate:"required" json:"username"`
	Password string `gorm:"type:varchar(255); not null" validate:"required"`
	Role     string `gorm:"type:varchar(255); not null" json:"role"`
	ModelDefault
}
