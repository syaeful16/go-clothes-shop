package models

type Cart struct {
	ID              uint    `gorm:"primaryKey" json:"id"`
	UserID          uint    `json:"user_id"`
	DetailProductID uint    `json:"product_id"`
	Quantity        uint    `gorm:"type:int; not null" json:"quantity" validate:"required,number"`
	TotalPrice      float32 `gorm:"type:decimal(10, 2); not null" json:"total_price" validate:"required,number"`
	ModelDefault
}
