package models

type UserAddress struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	RecipientName string `gorm:"type:varchar(255); not null" json:"recipient_name" validate:"required,alpha"`
	PhoneNumber   string `gorm:"type:varchar(255); not null" json:"phone_number" validate:"required,number"`
	Address       string `gorm:"type:text; not null" json:"address" validate:"required"`
	DetailAddress string `gorm:"type:varchar(255); not null" json:"detail_address" validate:"required,min=25"`
	AddressName   string `gorm:"type:varchar(255); not null" json:"address_name" validate:"required"`
	Status        bool   `gorm:"type:bool" json:"status"`
	UserID        uint   `json:"user_id"`
}
