package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = "secret"
var hmacSecretKey = []byte(secretKey)

var signingMethod = jwt.SigningMethodHS256

// Generates a JWT for a given username
func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24).Unix()
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

// Parses and validates a JWT string, returning the token if valid, else an error
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
	
