package models

import "github.com/dgrijalva/jwt-go"

type DepClaims struct {
	jwt.StandardClaims
}
