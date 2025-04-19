package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = "secret"
var hmacSecretKey = []byte(secretKey)

var signingMethod = jwt.SigningMethodHS256

// In a real api, I could replace this with an identifier
// or host of the api. More important with microservice architecture,
// but doesn't hurt with even one
var expectedIssuer = "api.mydomain.com"
var expectedAudience = expectedIssuer

// Generates a JWT for a given user
func GenerateToken(userid string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 1).Unix()
	currentTime := time.Now().Unix()
	notBefore := time.Now().Add(time.Millisecond * 100).Unix()
	token := jwt.NewWithClaims(signingMethod, jwt.MapClaims{
		"sub": userid,
		"iss": expectedIssuer,
		"aud": expectedAudience,
		"exp": expirationTime,
		"iat": currentTime,
		"nbf": notBefore,
	})

	return token.SignedString([]byte(secretKey))
}

// Parses and validates a JWT string, returning the token if valid, else an error
func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (any, error) {
			return hmacSecretKey, nil
		},
		jwt.WithValidMethods([]string{signingMethod.Alg()}),
		jwt.WithAudience(expectedAudience),
		jwt.WithIssuer(expectedIssuer),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	} else {
		return token, nil
	}
}
