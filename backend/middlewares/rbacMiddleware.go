package middlewares

import (
	"canny-clone/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// RoleRequired checks if the user has the required role
func RoleRequired(roles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Remove "Bearer " prefix if present
			if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
				tokenString = tokenString[7:]
			}

			claims := &services.TokenClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(services.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if user role is in the list of allowed roles
			roleMatched := false
			for _, role := range roles {
				if claims.Role == role {
					roleMatched = true
					break
				}
			}

			if !roleMatched {
				http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
				return
			}

			// Add user info to request for later use
			r.Header.Set("User-ID", strconv.Itoa(claims.UserID))
			r.Header.Set("User-Role", claims.Role)

			next(w, r)
		}
	}
}

// BoardRoleRequired checks if the user has the required role for a specific board
func BoardRoleRequired(boardRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// First ensure the user is authenticated
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
				tokenString = tokenString[7:]
			}

			claims := &services.TokenClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(services.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// If user is app_admin, they have access to all boards with all permissions
			if claims.Role == "app_admin" {
				r.Header.Set("User-ID", strconv.Itoa(claims.UserID))
				r.Header.Set("User-Role", claims.Role)
				next(w, r)
				return
			}

			// Get board ID from URL
			parts := strings.Split(r.URL.Path, "/")
			var boardIDStr string
			for i, part := range parts {
				if part == "boards" && i+1 < len(parts) {
					boardIDStr = parts[i+1]
					break
				}
			}

			if boardIDStr == "" {
				http.Error(w, "Invalid board ID", http.StatusBadRequest)
				return
			}

			boardID, err := strconv.Atoi(boardIDStr)
			if err != nil {
				http.Error(w, "Invalid board ID", http.StatusBadRequest)
				return
			}

			// Check user's role for this board
			userRepo := services.GetUserRepository()
			boardRoles, err := userRepo.GetUserBoardRoles(claims.UserID)
			if err != nil {
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}

			userBoardRole, exists := boardRoles[boardID]
			if !exists {
				http.Error(w, "Forbidden: Not a member of this board", http.StatusForbidden)
				return
			}

			// Check if the user's board role is sufficient
			roleMatched := false
			for _, role := range boardRoles {
				if userBoardRole == role {
					roleMatched = true
					break
				}
			}

			if !roleMatched {
				http.Error(w, "Forbidden: Insufficient permissions for this board", http.StatusForbidden)
				return
			}

			// Add user and board info to request for later use
			r.Header.Set("User-ID", strconv.Itoa(claims.UserID))
			r.Header.Set("User-Role", claims.Role)
			r.Header.Set("Board-Role", userBoardRole)
			
			next(w, r)
		}
	}
}

// GetUserFromRequest extracts user information from JWT token in request
func GetUserFromRequest(r *http.Request) (int, string, error) {
	userIDStr := r.Header.Get("User-ID")
	userRole := r.Header.Get("User-Role")
	
	if userIDStr == "" || userRole == "" {
		return 0, "", fmt.Errorf("user information not found in request")
	}
	
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, "", fmt.Errorf("invalid user ID in request")
	}
	
	return userID, userRole, nil
}
