package main

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
)

// TestAuthFlow tests the complete authentication flow:
// 1. Register a user
// 2. Login with those credentials
// 3. Use the token to access a protected endpoint
// Preconditions: db is reset, api is running. See Makefile, "make dev" starts it up, "make test" runs test with reset db
func TestAuthFlow(t *testing.T) {
	// Create httpexpect instance
	expect := httpexpect.Default(t, "http://localhost:8080")

	// Test user credentials
	email := "test@example.com"
	password := "testpassword123"

	// 1. Test Registration
	expect.POST("/register").
		WithJSON(map[string]string{
			"email":    email,
			"password": password,
		}).
		Expect().
		Status(http.StatusOK)

	// 2. Test Login
	cookie := expect.POST("/login").
		WithJSON(map[string]string{
			"email":    email,
			"password": password,
		}).
		Expect().
		Status(http.StatusOK).
		Cookie("token")

	// Get the token from cookies
	token := cookie.Value().Raw();
	assert.NotEmpty(t, token, "Expected token cookie in login response")

	// 3. Test Protected Endpoint
	expect.GET("/users").
		WithCookie("token", token).
		Expect().
		Status(http.StatusOK)
}

// TestInvalidToken verifies that protected endpoints reject invalid tokens
func TestInvalidToken(t *testing.T) {
	expect := httpexpect.Default(t, "http://localhost:8080")

	// Try to access /users with an invalid token
	expect.GET("/users").
		WithCookie("token", "invalid-token").
		Expect().
		Status(http.StatusUnauthorized)
}