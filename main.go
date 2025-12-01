// main.go
package main

import (
	"log" // Sekarang library ini akan terpakai di bawah
	"net/http"
	"os"
	"strings"

	"go-mongo-geojson/config"
	"go-mongo-geojson/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// Middleware Auth (Pengecekan Token)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			return
		}

		// Hapus "Bearer " dari string token jika ada
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "rahasia_super_aman"
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Next()
	}
}

func main() {
	// 1. Load .env (Local only)
	// Kita gunakan log di sini agar import "log" valid
	if err := godotenv.Load(); err != nil {
		log.Println("Info: .env file not found, using system environment variables")
	}

	// 2. Connect Database
	mongoClient := config.ConnectDB()

	// 3. Init Collections
	controllers.InitPlaceCollection(mongoClient)
	controllers.InitUserCollection(mongoClient) 

	// 4. Setup Router
	router := gin.Default()
	router.Use(cors.Default())

	// Static Files (Frontend)
	router.Static("/public", "./public")
	router.StaticFile("/", "./public/index.html")

	// 5. Routes API
	api := router.Group("/api")
	{
		// Public Routes
		api.GET("/places", controllers.GetPlaces) 
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)

		// Protected Routes (Butuh Login)
		protected := api.Group("/")
		protected.Use(AuthMiddleware()) 
		{
			protected.POST("/places", controllers.CreatePlace)
			protected.PUT("/places/:id", controllers.UpdatePlace)
			protected.DELETE("/places/:id", controllers.DeletePlace)
		}
	}

	// 6. Run Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Kita gunakan log di sini juga sebagai indikator server jalan
	log.Printf("Server running on port %s", port)
	router.Run(":" + port)
}