package cartcontroller

import (
	"encoding/json"
	"fmt"
	"go-clothes-shop/helper"
	"go-clothes-shop/middlewares"
	"go-clothes-shop/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type cart struct {
	ID            uint          `json:"id"`
	Quantity      uint          `json:"quantity"`
	TotalPrice    float32       `json:"total_price"`
	DetailProduct detailProduct `json:"detail_product"`
}

type detailProduct struct {
	ID      uint    `json:"id"`
	Photo   string  `json:"photo"`
	Color   string  `json:"color"`
	Size    string  `json:"size"`
	Stock   uint    `json:"stock"`
	Price   float32 `json:"price"`
	Product product `json:"product"`
}

type product struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Material string `json:"material"`
	Category string `json:"category"`
}

func Index(w http.ResponseWriter, r *http.Request) {

}

func Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	var cartExist models.Cart
	if err := models.DB.First(&cartExist, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := map[string]string{"message": err.Error()}
			helper.JSONResponse(w, http.StatusNotFound, response)
			return
		}
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	var detailProductExist models.DetailProduct
	if err := models.DB.First(&detailProductExist, cartExist.DetailProductID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := map[string]string{"message": err.Error()}
			helper.JSONResponse(w, http.StatusNotFound, response)
			return
		}
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	var productExist models.Product
	if err := models.DB.Where("id_product = ?", detailProductExist.ProductId).First(&productExist).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := map[string]string{"message": err.Error()}
			helper.JSONResponse(w, http.StatusNotFound, response)
			return
		}
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	result := cart{
		ID:         cartExist.ID,
		Quantity:   cartExist.Quantity,
		TotalPrice: cartExist.TotalPrice,
		DetailProduct: detailProduct{
			ID:    detailProductExist.ID,
			Photo: detailProductExist.Photo,
			Color: detailProductExist.Color,
			Size:  detailProductExist.Size,
			Price: detailProductExist.Price,
			Product: product{
				ID:       productExist.ID,
				Name:     productExist.Name,
				Material: productExist.Material,
				Category: productExist.Category,
			},
		},
	}

	helper.JSONResponse(w, http.StatusOK, result)
}

func Store(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IdKey).(uint)

	// get data from request body
	var cartInput models.Cart
	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&cartInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	var status bool = false

	// cek apakah product (detail product) sudah ditambahkan di cart atau belum
	var cartsExist models.Cart
	if err := models.DB.Where("detail_product_id = ? AND user_id = ?", cartInput.DetailProductID, userID).First(&cartsExist).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			status = true
		}

		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	fmt.Println("staus :", status)

	if !status {
		response := map[string]string{"message": "Product is already in the cart"}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	// get data detail product
	var detailProduct models.DetailProduct
	if err := models.DB.First(&detailProduct, cartInput.DetailProductID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := map[string]string{"message": err.Error() + " product"}
			helper.JSONResponse(w, http.StatusNotFound, response)
			return
		}
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	// check quantity cart inputted is greater than stock
	if detailProduct.Stock < cartInput.Quantity {
		response := map[string]string{"message": "The quantity you entered exceeds stock"}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	// jika quantity lebih kecil dari stok
	// maka kalkulasi sub total
	subTotal := detailProduct.Price * float32(cartInput.Quantity)

	// tampung dalam struct cart
	cartInput = models.Cart{
		UserID:          userID,
		DetailProductID: detailProduct.ID,
		Quantity:        cartInput.Quantity,
		TotalPrice:      subTotal,
	}

	// cek validasi
	// add validator message response
	if validate := helper.Validation(cartInput); validate != nil {
		response := map[string]interface{}{"error": validate}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// jika validasi aman maka store ke database
	if models.DB.Create(&cartInput).RowsAffected == 0 {
		response := map[string]string{"message": "Failed to add cart"}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	// tampilkan respon berhasil
	response := map[string]string{"message": "Success add cart"}
	helper.JSONResponse(w, http.StatusCreated, response)
}

func Update(w http.ResponseWriter, r *http.Request) {
	// get id parameter
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	var cartInput models.Cart
	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&cartInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// add validator message response
	if validate := helper.Validation(cartInput); validate != nil {
		response := map[string]interface{}{"error": validate}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	if models.DB.Where("id = ?", id).Updates(&cartInput).RowsAffected == 0 {
		response := map[string]string{"message": "Failed to update this cart"}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "success update data cart"}
	helper.JSONResponse(w, http.StatusOK, response)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	var cart models.Cart
	if models.DB.Delete(&cart, id).RowsAffected == 0 {
		response := map[string]string{"message": "failed to delete this cart"}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Success for delete cart"}
	helper.JSONResponse(w, http.StatusOK, response)
}
