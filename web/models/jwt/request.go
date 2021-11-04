package jwt

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	jwt.StandardClaims
	Role           int
	OrganizationID int
	DepartmentID   int
}
