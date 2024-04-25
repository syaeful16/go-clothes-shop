package detailproductcontroller

import (
	"fmt"
	"go-clothes-shop/helper"
	"go-clothes-shop/models"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {

}

func Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["product_id"]
	if productID == "" {
		response := map[string]string{"message": "Parameter not valid or empty"}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	var detailProduct []models.DetailProduct
	if err := models.DB.Where("product_id = ?", productID).Find(&detailProduct).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	if len(detailProduct) == 0 {
		response := map[string]string{"message": "record not found"}
		helper.JSONResponse(w, http.StatusNotFound, response)
		return
	}

	helper.JSONResponse(w, http.StatusOK, detailProduct)
}

func Store(w http.ResponseWriter, r *http.Request) {

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

	color := r.FormValue("color")
	size := r.FormValue("size")
	stock, err := strconv.ParseUint(r.FormValue("stock"), 10, 64)
	if stock != 0 {
		if err != nil {
			response := map[string]string{"message": err.Error()}
			helper.JSONResponse(w, http.StatusBadRequest, response)
			return
		}
	}
	price, err := strconv.ParseFloat(r.FormValue("price"), 32)
	if price != 0 {
		if err != nil {
			response := map[string]string{"message": err.Error()}
			helper.JSONResponse(w, http.StatusBadRequest, response)
			return
		}
	}
	productID := r.FormValue("product_id")

	detailProductInput := models.DetailProduct{
		Photo:     randomFilename,
		Color:     color,
		Size:      size,
		Stock:     uint(stock),
		Price:     float32(price),
		ProductId: productID,
	}

	// add validation
	if validate := helper.Validation(detailProductInput); validate != nil {
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
	fmt.Println(destinationFile.Name())

	if models.DB.Create(&detailProductInput).RowsAffected == 0 {
		response := map[string]string{"message": "Failed to add detail product"}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	_, err = io.Copy(destinationFile, photo)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Success add detail product"}
	helper.JSONResponse(w, http.StatusCreated, response)
}

func Update(w http.ResponseWriter, r *http.Request) {
	// get id in parameter
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

	// get data existing
	var detailProductExisting models.DetailProduct
	if err := models.DB.First(&detailProductExisting, id).Error; err != nil {
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

	color := r.FormValue("color")
	size := r.FormValue("size")
	stock, err := strconv.ParseUint(r.FormValue("stock"), 10, 64)
	if stock != 0 {
		if err != nil {
			response := map[string]string{"message": err.Error()}
			helper.JSONResponse(w, http.StatusBadRequest, response)
			return
		}
	}
	price, err := strconv.ParseFloat(r.FormValue("price"), 32)
	if price != 0 {
		if err != nil {
			response := map[string]string{"message": err.Error()}
			helper.JSONResponse(w, http.StatusBadRequest, response)
			return
		}
	}

	detailProductUpdate := models.DetailProduct{
		Photo:     randomFilename,
		Color:     color,
		Size:      size,
		Stock:     uint(stock),
		Price:     float32(price),
		ProductId: detailProductExisting.ProductId,
	}

	// add validation
	if validate := helper.Validation(detailProductUpdate); validate != nil {
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
	fmt.Println(destinationFile.Name())

	if models.DB.Where("id = ?", id).Updates(&detailProductUpdate).RowsAffected == 0 {
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

	// make filepath image
	filePath := filepath.Join("uploads", detailProductExisting.Photo)

	// Delete file
	errRemove := os.Remove(filePath)
	if errRemove != nil {
		fmt.Println("Error deleting file:", err)
	}

	fmt.Println("File deleted successfully.")

	response := map[string]string{"message": "Success for update detail product"}
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

	var detailProduct models.DetailProduct
	if models.DB.Delete(&detailProduct, id).RowsAffected == 0 {
		response := map[string]string{"message": "failed to delete this product"}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Success for delete product"}
	helper.JSONResponse(w, http.StatusOK, response)
}
