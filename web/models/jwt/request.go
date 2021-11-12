package jwt

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	jwt.StandardClaims
	Role           int
	OrganizationID int
	DepartmentID   int
}

type WXClaims struct {
	jwt.StandardClaims
	UID        int32
	Role       int
	Openid     string
	SessionKey string
}
