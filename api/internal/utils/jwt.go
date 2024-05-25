package utils

import (
	"errors"
	"time"

	"github.com/palp1tate/FlowFederate/api/internal/errorx"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	ID   int
	Role int
	jwt.StandardClaims
}

type JWT struct {
	SigningKey []byte
	Expiration time.Duration
}

func NewJWT(signingKey string, expiration int) *JWT {
	return &JWT{
		SigningKey: []byte(signingKey),
		Expiration: time.Duration(expiration) * time.Minute,
	}
}

func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Add(j.Expiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(j.SigningKey)
	if err != nil {
		return "", err
	}
	return "Bearer " + t, nil
}

func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
		return nil, errorx.TokenInvalid
	}
	tokenString = tokenString[7:]
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errorx.TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errorx.TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errorx.TokenNotValidYet
			} else {
				return nil, errorx.TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, errorx.TokenInvalid
}

func (j *JWT) RefreshToken(tokenString string) (string, error) {
	if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
		return "", errorx.TokenInvalid
	}
	tokenString = tokenString[7:]
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		claims.StandardClaims.ExpiresAt = time.Now().Add(j.Expiration).Unix()
		return j.CreateToken(*claims)
	}
	return "", errorx.TokenInvalid
}
