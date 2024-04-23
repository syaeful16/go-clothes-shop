package middlewares

import (
	"context"
	"fmt"
	"go-clothes-shop/config"
	"go-clothes-shop/helper"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const (
	IdKey       ContextKey = "id"
	UsernameKey ContextKey = "username"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				response := map[string]string{"message": "Unauthorized"}
				helper.JSONResponse(w, http.StatusUnauthorized, response)
				return
			}
		}
		// mengambil token value
		tokenString := c.Value

		claims := &config.JWTClaims{}
		// parsing token jwt
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		if err != nil {
			response := map[string]string{"message": "Unauthorized"}
			helper.JSONResponse(w, http.StatusUnauthorized, response)
			return
		}

		if !token.Valid {
			response := map[string]string{"message": "Unauthorized"}
			helper.JSONResponse(w, http.StatusUnauthorized, response)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AdminRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				response := map[string]string{"message": "Unauthorized"}
				helper.JSONResponse(w, http.StatusUnauthorized, response)
				return
			}
		}
		// mengambil token value
		tokenString := c.Value

		claims := &config.JWTClaims{}
		// parsing token jwt
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		if err != nil {
			response := map[string]string{"message": "Unauthorized"}
			helper.JSONResponse(w, http.StatusUnauthorized, response)
			return
		}

		if !token.Valid {
			response := map[string]string{"message": "Unauthorized"}
			helper.JSONResponse(w, http.StatusUnauthorized, response)
			return
		}

		fmt.Println(claims.UserId)
		fmt.Println(claims.Role)

		if claims.Role != "admin" {
			response := map[string]string{"message": "Your account does not have access"}
			helper.JSONResponse(w, http.StatusForbidden, response)
			return
		}

		ctx := context.WithValue(r.Context(), "id", claims.UserId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
