package model

import "github.com/golang-jwt/jwt/v4"

type UserClaims struct {
	jwt.RegisteredClaims
	Name string `json:"name"`
	Role string `json:"role"`
}
