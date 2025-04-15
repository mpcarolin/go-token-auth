package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = "secret"
var hmacSecretKey = []byte(secretKey)

var signingMethod = jwt.SigningMethodHS256

func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24)
	currentTime := time.Now().Unix()
	notBefore := time.Now().Add(time.Millisecond * 100).Unix()
	token := jwt.NewWithClaims(signingMethod, jwt.MapClaims{
		"username": username,
		"exp":      expirationTime,
		"iat":      currentTime,
		"nbf":      notBefore,
	})

	return token.SignedString([]byte(secretKey))
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return hmacSecretKey, nil
	}, jwt.WithValidMethods([]string{signingMethod.Alg()}))

	if err != nil {
		return nil, err
	} else {
		return token, nil
	}
}
	
