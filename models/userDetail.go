package models

type UserDetail struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Fullname    string `gorm:"type:varchar(255); not null" json:"fullname" validate:"required"`
	PhoneNumber string `gorm:"type:varchar(255); not null" json:"phone_number" validate:"required,number"`
	Email       string `gorm:"type:varchar(255)" json:"email" validate:"email"`
	Gender      string `gorm:"type:varchar(255)" json:"gender" validate:"required"`
	UserID      uint   `json:"user_id"`
	User        User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
