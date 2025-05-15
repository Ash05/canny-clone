package services

import (
	"canny-clone/repositories"
	"canny-clone/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleOAuthConfig holds OAuth configuration
var (
	GoogleOAuthConfig *oauth2.Config
	JWTSecret         string
)

// GoogleUserInfo represents the user info returned from Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// TokenClaims represents the JWT claims
type TokenClaims struct {
	UserID int    `json:"userId"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Role   string `json:"role"` // Global role: "app_admin", "stakeholder", or "user"
	jwt.StandardClaims
}

// InitAuth initializes authentication configuration
func InitAuth() {
	config := utils.GetConfig()
	
	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     config.GoogleClientID,
		ClientSecret: config.GoogleClientSecret,
		RedirectURL:  config.GoogleRedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	
	JWTSecret = config.JWTSecret
}

// GetGoogleLoginURL generates the Google OAuth login URL
func GetGoogleLoginURL(w http.ResponseWriter, r *http.Request) {
	url := GoogleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": url})
}

// HandleGoogleCallback processes the OAuth callback from Google
func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Get the code from the query
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// Exchange the code for a token
	token, err := GoogleOAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// Get user info from Google
	client := GoogleOAuthConfig.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	userData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(userData, &userInfo); err != nil {
		http.Error(w, "Failed to parse user data", http.StatusInternalServerError)
		return
	}

	// Find or create the user in the database
	userRepo := repositories.NewUserRepository()
	user, err := userRepo.FindUserByEmail(userInfo.Email)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		// Create new user
		newUser := &repositories.User{
			Email:    userInfo.Email,
			Name:     userInfo.Name,
			Picture:  userInfo.Picture,
			Provider: "google",
		}
		
		if err := userRepo.CreateUser(newUser); err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		
		// Get the created user
		user, err = userRepo.FindUserByEmail(userInfo.Email)
		if err != nil {
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
			return
		}
	}

	// Generate JWT token
	jwtToken, err := generateJWT(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return the token to the client
	// In a real app, you might want to redirect to the frontend with the token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": jwtToken,
		"name":  user.Name,
		"email": user.Email,
	})
}

// Generate JWT token for a user
func generateJWT(user *repositories.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	
	// Get user's board roles if needed for enhanced token
	userRepo := repositories.NewUserRepository()
	boardRoles, err := userRepo.GetUserBoardRoles(user.ID)
	if err != nil {
		return "", err
	}
	
	// Convert board roles to JSON string
	boardRolesJSON, err := json.Marshal(boardRoles)
	if err != nil {
		return "", err
	}
	
	claims := &TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// AuthMiddleware checks if the request has a valid JWT token
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
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

		claims := &TokenClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user ID to request context or headers for downstream handlers
		r.Header.Set("User-ID", string(claims.UserID))
		next(w, r)
	}
}

// GetUserRepository returns an instance of the user repository
func GetUserRepository() repositories.UserRepository {
	return repositories.NewUserRepository()
}
