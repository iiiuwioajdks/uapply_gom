package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"uapply_go/web/global"
	jwt2 "uapply_go/web/models/jwt"
)

type JWT struct {
	SigningKey []byte
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("uapply-token")
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
		c.Set("claim", claim)
		zap.S().Info(claim)
		c.Next()
	}
}

// WXJWTAuth 微信小程序生成token
func WXJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("uapply-token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "无携带token，请先登录",
			})
			c.Abort()
			return
		}

		j := new(JWT)
		claim, err := j.ParseWXToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": err.Error(),
			})
			c.Abort()
			return
		}
		// set 一些东西
		c.Set("wxClaim", claim)
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
func (j *JWT) CreateToken(claims jwt2.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// CreateWXToken 生成jwt
func (j *JWT) CreateWXToken(claims jwt2.WXClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

var (
	TokenMalformed   = errors.New("That is not even a token")
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenInvalid     = errors.New("Couldn't handle this token")
)

func (j *JWT) ParseToken(tokenString string) (*jwt2.Claims, error) {
	var cs jwt2.Claims
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

func (j *JWT) ParseWXToken(tokenString string) (*jwt2.WXClaims, error) {
	var cs jwt2.WXClaims
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
