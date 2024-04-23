package productcontroller

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

type productResponse struct {
	ID          uint         `json:"id"`
	IdProduct   string       `json:"id_product"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Material    string       `json:"material"`
	Category    string       `json:"category"`
	UserId      userResponse `json:"created_by"`
	CreatedAt   string       `json:"created_at"`
	UpdatedAt   string       `json:"updated_at"`
}

type userResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"name"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	if err := models.DB.Find(&products).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	if len(products) == 0 {
		response := map[string]string{"message": "record not found"}
		helper.JSONResponse(w, http.StatusNotFound, response)
		return
	}

	responses := make([]productResponse, len(products))
	for i, product := range products {
		var user models.User
		if err := models.DB.Where("id = ?", products[i].UserID).First(&user).Error; err != nil {
			response := map[string]string{"message": err.Error()}
			helper.JSONResponse(w, http.StatusInternalServerError, response)
			return
		}

		responses[i] = productResponse{
			ID:          product.ID,
			IdProduct:   product.IdProduct,
			Name:        product.Name,
			Description: product.Description,
			Material:    product.Material,
			Category:    product.Category,
			UserId: userResponse{
				ID:       user.ID,
				Username: user.Username,
			},
			CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	helper.JSONResponse(w, http.StatusOK, responses)

}

func Show(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	if err := models.DB.Where("id = ?", id).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := map[string]string{"message": err.Error()}
			helper.JSONResponse(w, http.StatusNotFound, response)
			return
		}
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	var user models.User
	if err := models.DB.Where("id = ?", product.UserID).First(&user).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	productResponse := productResponse{
		ID:          product.ID,
		IdProduct:   product.IdProduct,
		Name:        product.Name,
		Description: product.Description,
		Material:    product.Material,
		Category:    product.Category,
		UserId: userResponse{
			ID:       user.ID,
			Username: user.Username,
		},
		CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	helper.JSONResponse(w, http.StatusOK, productResponse)
}

func Store(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IdKey).(uint)
	fmt.Println(userID)

	var productInput models.Product
	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&productInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// add validator message response
	if validate := helper.Validation(productInput); validate != nil {
		response := map[string]interface{}{"error": validate}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	productInput.UserID = userID
	if models.DB.Create(&productInput).RowsAffected == 0 {
		response := map[string]string{"message": "Failed to add product"}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Success add product"}
	helper.JSONResponse(w, http.StatusOK, response)
}

func Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	var productInput models.Product
	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&productInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// add validator message response
	if validate := helper.Validation(productInput); validate != nil {
		response := map[string]interface{}{"error": validate}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	if models.DB.Where("id = ?", id).Updates(&productInput).RowsAffected == 0 {
		response := map[string]string{"message": "Failed to update this product"}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "success update data product"}
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

	var product models.Product
	if models.DB.Delete(&product, id).RowsAffected == 0 {
		response := map[string]string{"message": "failed to delete this product"}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Success for delete product"}
	helper.JSONResponse(w, http.StatusOK, response)
}
