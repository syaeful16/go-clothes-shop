package models

type Product struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	IdProduct     string          `gorm:"type:varchar(255); not null; unique" json:"id_product" validate:"required"`
	Name          string          `gorm:"type:varchar(255); not null" json:"name" validate:"required"`
	Description   string          `gorm:"type:varchar(255); not null" json:"description" validate:"required,min=30"`
	Material      string          `gorm:"type:varchar(255); not null" json:"material" validate:"required"`
	Category      string          `gorm:"type:varchar(255); not null" json:"category" validate:"required"`
	UserID        uint            `json:"user_id"`
	DetailProduct []DetailProduct `gorm:"foreignKey:ProductId; references:IdProduct; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ModelDefault
}
