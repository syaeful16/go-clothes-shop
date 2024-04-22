package config

import "github.com/golang-jwt/jwt/v5"

var JWT_KEY = []byte("c9e5cc2c7cfcfe4c7b94f28a695deeb2b67daa92d2005b1b35d760b138927657")

type JWTClaims struct {
	UserId uint
	Role   string
	jwt.RegisteredClaims
}
