package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"uapply_go/web/global"
	"uapply_go/web/models"
)

type JWT struct {
	SigningKey []byte
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "无携带token",
			})
			c.Abort()
			return
		}

		j := new(JWT)
		claim, err := j.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": err.Error(),
			})
			c.Abort()
			return
		}
		// set 一些东西
		zap.S().Info(claim)
		c.Next()
	}
}

func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte(global.Conf.JwtInfo.SigningKey),
	}
}

// CreateToken 生成jwt
func (j *JWT) CreateToken(claims models.DepClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

var (
	TokenMalformed   = errors.New("That is not even a token")
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenInvalid     = errors.New("Couldn't handle this token")
)

func (j *JWT) ParseToken(tokenString string) (*models.DepClaims, error) {
	var cs models.DepClaims
	_, err := jwt.ParseWithClaims(tokenString, &cs, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Conf.JwtInfo.SigningKey), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		} else {
			return nil, TokenInvalid
		}
	}
	return &cs, nil
}
