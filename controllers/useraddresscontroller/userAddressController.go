package useraddresscontroller

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

func Index(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IdKey).(uint)

	var userAddressExist []models.UserAddress
	if err := models.DB.Where("user_id = ?", userID).Find(&userAddressExist).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	if len(userAddressExist) == 0 {
		response := map[string]string{"message": "record not found address"}
		helper.JSONResponse(w, http.StatusNotFound, response)
		return
	}

	helper.JSONResponse(w, http.StatusOK, userAddressExist)
}

func Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	var userAddressExist models.UserAddress
	if err := models.DB.First(&userAddressExist, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := map[string]string{"message": err.Error() + " address"}
			helper.JSONResponse(w, http.StatusNotFound, response)
			return
		}
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	helper.JSONResponse(w, http.StatusOK, userAddressExist)
}

func Store(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IdKey).(uint)

	var userAddressInput models.UserAddress
	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&userAddressInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	userAddressInput.UserID = userID
	if validate := helper.Validation(userAddressInput); validate != nil {
		response := map[string]interface{}{"message": validate}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	if models.DB.Create(&userAddressInput).RowsAffected == 0 {
		response := map[string]string{"message": "Failed to add address"}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	// tampilkan respon berhasil
	response := map[string]string{"message": "Success add address"}
	helper.JSONResponse(w, http.StatusCreated, response)
}

func Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IdKey).(uint)

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	var userAddressInput models.UserAddress
	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&userAddressInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	userAddressInput.UserID = userID
	// add validator message response
	if validate := helper.Validation(userAddressInput); validate != nil {
		response := map[string]interface{}{"error": validate}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	fmt.Println(userAddressInput.Address, userAddressInput.AddressName, userAddressInput.DetailAddress, userAddressInput.PhoneNumber, userAddressInput.RecipientName)
	if models.DB.Where("id = ?", id).Updates(&userAddressInput).RowsAffected == 0 {
		response := map[string]string{"message": "Failed to update this address"}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "success update addres"}
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

	var userAddressExist models.UserAddress
	if models.DB.Delete(&userAddressExist, id).RowsAffected == 0 {
		response := map[string]string{"message": "failed to delete this address"}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Success for delete address"}
	helper.JSONResponse(w, http.StatusOK, response)
}
