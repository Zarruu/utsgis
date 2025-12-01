// main.go
package main

import (
	"log"
	"os" // Tambahkan import os untuk membaca env var

	"go-mongo-geojson/config"
	"go-mongo-geojson/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Muat file .env (Hanya untuk Local Development)
	// Di Cloud (Render/Koyeb), file .env tidak ikut di-upload.
	// Variabel akan dibaca langsung dari sistem environment hosting.
	err := godotenv.Load()
	if err != nil {
		log.Println("Info: .env file not found, using system environment variables")
	}

	// 2. Hubungkan ke database
	// Pastikan fungsi ini membaca os.Getenv("MONGO_URI")!
	mongoClient := config.ConnectDB()

	// PANGGIL FUNGSI INISIALISASI CONTROLLER DI SINI
	controllers.InitPlaceCollection(mongoClient)

	// Inisialisasi router Gin
	// Di production, set GIN_MODE=release di environment variable hosting
	router := gin.Default()

	// Setup CORS
	// Default() mengizinkan semua origin, aman untuk tahap awal/testing.
	// Nanti bisa diperketat jika frontend sudah punya domain tetap.
	router.Use(cors.Default())

	// Definisikan rute API
	api := router.Group("/api")
	{
		api.POST("/places", controllers.CreatePlace)
		api.GET("/places", controllers.GetPlaces)
		api.PUT("/places/:id", controllers.UpdatePlace)
		api.DELETE("/places/:id", controllers.DeletePlace)
	}

	// 3. Konfigurasi PORT Dinamis (WAJIB UNTUK CLOUD)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback ke 8080 jika dijalankan di local tanpa .env
	}

	// Jalankan server dengan port dinamis
	log.Printf("Server running on port %s", port)
	router.Run(":" + port) 
}