package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = os.Getenv("TOKEN_SECRET_KEY");
var hmacSecretKey = []byte(secretKey)

var signingMethod = jwt.SigningMethodHS256

// In a real api, I could replace this with an identifier
// or host of the api. More important with microservice architecture,
// but doesn't hurt with even one
var expectedIssuer = "api.mydomain.com"
var expectedAudience = expectedIssuer

// Returns the time (in minutes via time.Duration) that a user can stay
// authenticated after logging in. Useful for token AND cookie expiration
func GetAuthDuration() time.Duration {
	durationStr := os.Getenv("AUTH_DURATION_MINUTES")
	if durationStr == "" {
		return time.Hour // Default to 1 hour
	}
	
	duration, err := strconv.Atoi(durationStr)
	if err != nil || duration <= 0 {
		return time.Hour // Default to 1 hour if invalid
	}
	
	return time.Duration(duration) * time.Minute
}

// Generates a JWT for a given user
func GenerateToken(userid string) (string, error) {
	authDuration := GetAuthDuration()
	expirationTime := time.Now().Add(authDuration).Unix()
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
