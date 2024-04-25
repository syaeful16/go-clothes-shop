package models

type DetailProduct struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Photo     string  `gorm:"type:varchar(255); not null" json:"photo" validate:"required"`
	Color     string  `gorm:"type:varchar(255); not null" json:"color" validate:"required"`
	Size      string  `gorm:"type:varchar(255); not null" json:"size" validate:"required"`
	Stock     uint    `gorm:"type:INT; not null" json:"stock" validate:"required,number"`
	Price     float32 `gorm:"type:decimal(10, 2); not null" json:"price" validate:"required,number"`
	ProductId string  `gorm:"type:varchar(255); not null" json:"product_id" validate:"required"`
	ModelDefault
}

/*
	{
  "photo": "image.jpg",
  "color": "blue",
  "size": "M",
  "stock": 2,
  "price": 10000.00,
  "product_id": "cloth1"
}
*/
