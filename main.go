// main.go
package main

import (
	"log"
	"os"

	"go-mongo-geojson/config"
	"go-mongo-geojson/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load .env (Local only)
	err := godotenv.Load()
	if err != nil {
		log.Println("Info: .env file not found, using system environment variables")
	}

	// 2. Connect Database
	mongoClient := config.ConnectDB()

	// 3. Init Collection
	controllers.InitPlaceCollection(mongoClient)

	// 4. Init Router
	router := gin.Default()

	// Setup CORS
	router.Use(cors.Default())

	// --- [PERBAIKAN UTAMA DI SINI] ---
	// Melayani file static (CSS/JS jika ada di dalam folder public)
	router.Static("/public", "./public")
	
	// Melayani file index.html di root URL "/"
	// Pastikan folder "public" dan file "index.html" ikut ter-upload ke Render!
	router.StaticFile("/", "./public/index.html") 
	// ----------------------------------

	// Group API
	api := router.Group("/api")
	{
		api.POST("/places", controllers.CreatePlace)
		api.GET("/places", controllers.GetPlaces)
		api.PUT("/places/:id", controllers.UpdatePlace)
		api.DELETE("/places/:id", controllers.DeletePlace)
	}

	// 5. Run Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	router.Run(":" + port)
}