package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
	"restapi-go/config"
	"restapi-go/models"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// TokenBlacklist stores invalidated tokens
var TokenBlacklist = make(map[string]time.Time)

// Login an existing user
func Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Fetch user from the database
	var user models.User
	err := config.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email=$1", credentials.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := generateJWT(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Respond with the token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"message": "Login successful",
	})
}

// Logout handles user logout by invalidating the JWT token
func Logout(w http.ResponseWriter, r *http.Request) {
	// Get the token from Authorization header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "No token provided", http.StatusBadRequest)
		return
	}

	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// Parse the token to get its expiration time
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.JwtSecret, nil
	})

	if err != nil {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	// Add token to blacklist
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		exp := time.Unix(int64(claims["exp"].(float64)), 0)
		TokenBlacklist[tokenString] = exp

		// Clean up expired tokens from blacklist
		go cleanupBlacklist()
	}

	// Clear any session cookies if you're using them
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour), // Set expiration in the past
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
		"status": "success",
	})
}

// Generate a JWT token
func generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["iat"] = time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JwtSecret)
}

// Get current authenticated user
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	var user models.User
	err := config.DB.QueryRow("SELECT id, name, email FROM users WHERE id=$1", userID).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Middleware to check if token is blacklisted
func IsTokenBlacklisted(tokenString string) bool {
	if exp, exists := TokenBlacklist[tokenString]; exists {
		if time.Now().Before(exp) {
			return true
		}
		// If token is expired, remove it from blacklist
		delete(TokenBlacklist, tokenString)
	}
	return false
}

// Clean up expired tokens from the blacklist
func cleanupBlacklist() {
	currentTime := time.Now()
	for token, exp := range TokenBlacklist {
		if currentTime.After(exp) {
			delete(TokenBlacklist, token)
		}
	}
}