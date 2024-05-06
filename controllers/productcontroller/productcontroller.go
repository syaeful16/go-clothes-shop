package productcontroller

import (
	"fmt"
	"go-clothes-shop/helper"
	"go-clothes-shop/middlewares"
	"go-clothes-shop/models"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type productResponse struct {
	ID          uint          `json:"id"`
	IdProduct   string        `json:"id_product"`
	Photo       string        `json:"photo_product"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Material    string        `json:"material"`
	Category    string        `json:"category"`
	Price       priceResponse `json:"price"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
}

type priceResponse struct {
	Min float32 `json:"min_price"`
	Max float32 `json:"max_price"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	if err := models.DB.Order("id DESC").Find(&products).Error; err != nil {
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
		var price priceResponse
		if err := models.DB.Model(&models.DetailProduct{}).Select("min(price) as min, max(price) as max").Where("product_id = ?", product.IdProduct).Scan(&price).Error; err != nil {
			response := map[string]string{"message": err.Error()}
			helper.JSONResponse(w, http.StatusInternalServerError, response)
			return
		}

		responses[i] = productResponse{
			ID:          product.ID,
			IdProduct:   product.IdProduct,
			Photo:       product.PhotoProduct,
			Name:        product.Name,
			Description: product.Description,
			Material:    product.Material,
			Category:    product.Category,
			Price: priceResponse{
				Min: price.Min,
				Max: price.Max,
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

	if err := models.DB.First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := map[string]string{"message": err.Error() + " product"}
			helper.JSONResponse(w, http.StatusNotFound, response)
			return
		}
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	var price priceResponse
	if err := models.DB.Model(&models.DetailProduct{}).Select("min(price) as min, max(price) as max").Where("product_id = ?", product.IdProduct).Scan(&price).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	productResponse := productResponse{
		ID:          product.ID,
		IdProduct:   product.IdProduct,
		Photo:       product.PhotoProduct,
		Name:        product.Name,
		Description: product.Description,
		Material:    product.Material,
		Category:    product.Category,
		Price: priceResponse{
			Min: price.Min,
			Max: price.Max,
		},
		CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	helper.JSONResponse(w, http.StatusOK, productResponse)
}

func Store(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IdKey).(uint)

	if err := r.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	photo, handler, err := r.FormFile("photo")
	if err != nil {
		// Check if the file is empty or not found
		if err == http.ErrMissingFile {
			// Handle the case where no file is uploaded
			response := map[string]string{"message": "No file uploaded"}
			helper.JSONResponse(w, http.StatusBadRequest, response)
			return
		}
		// Handle other errors
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}
	defer photo.Close()

	extValid := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	// Check if the file extension is allowed
	fileExt := strings.ToLower(filepath.Ext(handler.Filename))
	if !extValid[fileExt] {
		response := map[string]string{"message": "Extension is not supported"}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Check file size
	if handler.Size > (500 * 1024) {
		response := map[string]string{"message": "File size too large"}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Generate random filename
	randomFilename := uuid.New().String() + fileExt

	idProduct := r.FormValue("id_product")
	name := r.FormValue("name")
	description := r.FormValue("description")
	material := r.FormValue("material")
	category := r.FormValue("category")

	productInput := models.Product{
		IdProduct:    idProduct,
		PhotoProduct: randomFilename,
		Name:         name,
		Description:  description,
		Material:     material,
		Category:     category,
	}

	// add validator message response
	if validate := helper.Validation(productInput); validate != nil {
		response := map[string]interface{}{"error": validate}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Create destination file
	destinationFile, err := os.Create(filepath.Join("uploads", randomFilename))
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}
	defer destinationFile.Close()

	productInput.UserID = userID
	if models.DB.Create(&productInput).RowsAffected == 0 {
		response := map[string]string{"message": "Failed to add product"}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	// copy file ke destination
	_, err = io.Copy(destinationFile, photo)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Success add product"}
	helper.JSONResponse(w, http.StatusCreated, response)
}

func Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	if err := r.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	//get value id_product
	var productExist models.Product
	if err := models.DB.First(&productExist, id).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	photo, handler, err := r.FormFile("photo")
	if err != nil {
		// Check if the file is empty or not found
		if err == http.ErrMissingFile {
			// Handle the case where no file is uploaded
			response := map[string]string{"message": "No file uploaded"}
			helper.JSONResponse(w, http.StatusBadRequest, response)
			return
		}
		// Handle other errors
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}
	defer photo.Close()

	extValid := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	// Check if the file extension is allowed
	fileExt := strings.ToLower(filepath.Ext(handler.Filename))
	if !extValid[fileExt] {
		response := map[string]string{"message": "Extension is not supported"}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Check file size
	if handler.Size > (500 * 1024) {
		response := map[string]string{"message": "File size too large"}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Generate random filename
	randomFilename := uuid.New().String() + fileExt

	name := r.FormValue("name")
	description := r.FormValue("description")
	material := r.FormValue("material")
	category := r.FormValue("category")

	productInput := models.Product{
		IdProduct:    productExist.IdProduct,
		PhotoProduct: randomFilename,
		Name:         name,
		Description:  description,
		Material:     material,
		Category:     category,
	}

	// add validator message response
	if validate := helper.Validation(productInput); validate != nil {
		response := map[string]interface{}{"error": validate}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Create destination file
	destinationFile, err := os.Create(filepath.Join("uploads", randomFilename))
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}
	defer destinationFile.Close()

	if models.DB.Where("id = ?", id).Updates(&productInput).RowsAffected == 0 {
		response := map[string]string{"message": "Failed to update this product"}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	// copy file to directory
	_, err = io.Copy(destinationFile, photo)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	// make filepath image for delete
	filePath := filepath.Join("uploads", productExist.PhotoProduct)

	// Delete file
	errRemove := os.Remove(filePath)
	if errRemove != nil {
		fmt.Println("Error deleting file:", err)
		return
	}

	fmt.Println("File deleted successfully.")

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
