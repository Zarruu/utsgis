// main.go (Update bagian Route)
package main

import (
	// ... import lain tetap sama
	"log"
	"os"
	"strings"
	"net/http"
	
	"go-mongo-geojson/config"
	"go-mongo-geojson/controllers"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Middleware Sederhana untuk Cek Token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			return
		}

		// Hapus "Bearer " dari string token
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" { jwtSecret = "rahasia_super_aman" }

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
	godotenv.Load()
	mongoClient := config.ConnectDB()

	// INISIALISASI COLLECTION
	controllers.InitPlaceCollection(mongoClient)
	controllers.InitUserCollection(mongoClient) // <-- TAMBAHAN BARU

	router := gin.Default()
	router.Use(cors.Default())

	router.Static("/public", "./public")
	router.StaticFile("/", "./public/index.html")

	api := router.Group("/api")
	{
		// Public Routes (Siapapun bisa akses)
		api.GET("/places", controllers.GetPlaces) // Orang bisa lihat peta tanpa login
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)

		// Protected Routes (Harus Login)
		protected := api.Group("/")
		protected.Use(AuthMiddleware()) // Pasang Gembok di sini
		{
			protected.POST("/places", controllers.CreatePlace)
			protected.PUT("/places/:id", controllers.UpdatePlace)
			protected.DELETE("/places/:id", controllers.DeletePlace)
		}
	}

	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	router.Run(":" + port)
}