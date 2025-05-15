package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"canny-clone/routes"
	"canny-clone/services"
	"canny-clone/utils"
)

func main() {
	// Load configuration
	utils.LoadConfig()
	config := utils.GetConfig()
	
	// Initialize database connection
	services.InitDB()
	
	// Initialize authentication service
	services.InitAuth()

	// Create router and register routes
	r := mux.NewRouter()
	routes.RegisterBoardRoutes(r)
	routes.RegisterFeedbackRoutes(r)
	routes.RegisterCategoryRoutes(r)
	routes.RegisterCommentRoutes(r)
	routes.RegisterAuthRoutes(r)

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // In production, specify your frontend domain
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "User-ID"},
		AllowCredentials: true,
	})
	
	// Use CORS middleware
	handler := c.Handler(r)
	
	// Start server
	port := config.Port
	if port == "" {
		port = "8080"
	}
	
	fmt.Printf("Server is running on port %s\n", port)
	http.ListenAndServe(":"+port, handler)
}
