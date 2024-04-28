package userdetailcontroller

import (
	"encoding/json"
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

	var userDetailExist models.UserDetail
	if err := models.DB.Where("user_id = ?", userID).First(&userDetailExist).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := map[string]string{"message": err.Error() + " user detail"}
			helper.JSONResponse(w, http.StatusNotFound, response)
			return
		}
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	helper.JSONResponse(w, http.StatusOK, userDetailExist)
}

// func Show(w http.ResponseWriter, r *http.Request) {

// }

func Store(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IdKey).(uint)

	var userDetailInput models.UserDetail
	// if err := models.DB.Where("user_id = ?", userID).First(&userDetailInput).Error; err != nil {
	// 	if err == gorm.ErrRecordNotFound {
	// 		response := map[string]string{"message": err.Error()}
	// 		helper.JSONResponse(w, http.StatusNotFound, response)
	// 		return
	// 	}
	// 	response := map[string]string{"message": err.Error()}
	// 	helper.JSONResponse(w, http.StatusInternalServerError, response)
	// 	return
	// }

	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&userDetailInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	userDetailInput.UserID = userID
	if validate := helper.Validation(userDetailInput); validate != nil {
		response := map[string]interface{}{"message": validate}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	if models.DB.Create(&userDetailInput).RowsAffected == 0 {
		response := map[string]string{"message": "Failed to create user detail"}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	// tampilkan respon berhasil
	response := map[string]string{"message": "Success create user detail"}
	helper.JSONResponse(w, http.StatusCreated, response)
}

func Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IdKey).(uint)
	// vars := mux.Vars(r)
	// id, err := strconv.ParseInt(vars["id"], 10, 64)
	// if err != nil {
	// 	response := map[string]string{"message": err.Error()}
	// 	helper.JSONResponse(w, http.StatusBadRequest, response)
	// 	return
	// }

	var userDetailInput models.UserDetail
	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&userDetailInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// add validator message response
	if validate := helper.Validation(userDetailInput); validate != nil {
		response := map[string]interface{}{"error": validate}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	if models.DB.Where("user_id = ?", userID).Updates(&userDetailInput).RowsAffected == 0 {
		response := map[string]string{"message": "Failed to update user detail"}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "success update user detail"}
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

	var userDetailInput models.UserDetail
	if models.DB.Delete(&userDetailInput, id).RowsAffected == 0 {
		response := map[string]string{"message": "failed to delete this use detail"}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Success for delete user detail"}
	helper.JSONResponse(w, http.StatusOK, response)
}
