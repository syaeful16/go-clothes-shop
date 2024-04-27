package authcontroller

import (
	"encoding/json"
	"go-clothes-shop/config"
	"go-clothes-shop/helper"
	"go-clothes-shop/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var userInput models.User
	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// add validator message response
	if validate := helper.Validation(userInput); validate != nil {
		response := map[string]interface{}{"error": validate}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	var existingUser models.User
	if err := models.DB.Where("username = ?", userInput.Username).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := map[string]string{"message": "username or password invalid"}
			helper.JSONResponse(w, http.StatusNotFound, response)
			return
		}

		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(userInput.Password)); err != nil {
		response := map[string]string{"message": "username or password invalid"}
		helper.JSONResponse(w, http.StatusNotFound, response)
		return
	}

	expTime := time.Now().Add(time.Minute * 1)
	claims := &config.JWTClaims{
		Role:   existingUser.Role,
		UserId: existingUser.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-fashion-shop",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenAlgo.SignedString(config.JWT_KEY)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})

	response := map[string]string{
		"username": existingUser.Username,
		"role":     existingUser.Role,
	}

	helper.JSONResponse(w, http.StatusOK, response)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var userInput models.User
	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// add validator message response
	if validate := helper.Validation(userInput); validate != nil {
		response := map[string]interface{}{"error": validate}
		helper.JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	var existinguser models.User
	if err := models.DB.Where("username = ?", userInput.Username).First(&existinguser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
			userInput.Password = string(hashPassword)
			userInput.Role = "user"

			if errCreate := models.DB.Create(&userInput).Error; errCreate != nil {
				response := map[string]string{"message": err.Error()}
				helper.JSONResponse(w, http.StatusInternalServerError, response)
				return
			}

			response := map[string]string{"message": "Success register"}
			helper.JSONResponse(w, http.StatusCreated, response)
			return
		}

	}

	// username has already been created
	response := map[string]string{"message": "username has already been created"}
	helper.JSONResponse(w, http.StatusFound, response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	response := map[string]string{"message": "Successful Logout"}
	helper.JSONResponse(w, http.StatusOK, response)
}
