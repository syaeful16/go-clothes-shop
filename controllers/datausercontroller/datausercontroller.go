package datausercontroller

import (
	"go-clothes-shop/helper"
	"go-clothes-shop/middlewares"
	"go-clothes-shop/models"
	"net/http"
)

type cartResponse struct {
	ID            uint                  `json:"id"`
	Quantity      uint                  `json:"quantity"`
	TotalPrice    float32               `json:"total_price"`
	DetailProduct detailProductResponse `json:"detail_product"`
}

type detailProductResponse struct {
	ID      uint            `json:"id"`
	Photo   string          `json:"photo"`
	Color   string          `json:"color"`
	Size    string          `json:"size"`
	Price   float32         `json:"price"`
	Product productResponse `json:"product"`
}

type productResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Material string `json:"material"`
	Category string `json:"category"`
}

func ShowCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IdKey).(uint)

	var cartsExist []models.Cart
	if err := models.DB.Joins("JOIN users ON users.id = carts.user_id").Joins("JOIN detail_products ON detail_products.id = carts.detail_product_id").Where("users.id = ? AND detail_products.deleted_at is null", userID).Find(&cartsExist).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	if len(cartsExist) == 0 {
		response := map[string]string{"message": "record not found"}
		helper.JSONResponse(w, http.StatusNotFound, response)
		return
	}

	response := make([]cartResponse, len(cartsExist))
	for i, cart := range cartsExist {
		var detailProductExist models.DetailProduct
		if err := models.DB.First(&detailProductExist, cart.DetailProductID).Error; err != nil {
			response := map[string]string{"message": err.Error()}
			helper.JSONResponse(w, http.StatusInternalServerError, response)
		}

		var productExist models.Product
		if err := models.DB.Where("id_product = ?", detailProductExist.ProductId).First(&productExist).Error; err != nil {
			response := map[string]string{"message": err.Error()}
			helper.JSONResponse(w, http.StatusInternalServerError, response)
		}

		response[i] = cartResponse{
			ID:         cart.ID,
			Quantity:   cart.Quantity,
			TotalPrice: cart.TotalPrice,
			DetailProduct: detailProductResponse{
				ID:    detailProductExist.ID,
				Photo: detailProductExist.Photo,
				Color: detailProductExist.Color,
				Size:  detailProductExist.Size,
				Price: detailProductExist.Price,
				Product: productResponse{
					ID:       productExist.ID,
					Name:     productExist.Name,
					Material: productExist.Material,
					Category: productExist.Category,
				},
			},
		}
	}

	// fmt.Println()

	helper.JSONResponse(w, http.StatusOK, response)
}
