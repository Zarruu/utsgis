// controllers/auth_controller.go
package controllers

import (
	"context"
	"go-mongo-geojson/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection

// Init User Collection (Panggil ini di main.go)
func InitUserCollection(client *mongo.Client) {
	userCollection = client.Database("location_db").Collection("users")
}

// 1. REGISTER
func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash Password (Jangan simpan password asli!)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Cek apakah username sudah ada
	count, _ := userCollection.CountDocuments(ctx, bson.M{"username": user.Username})
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// 2. LOGIN
func Login(c *gin.Context) {
	var input models.User
	var foundUser models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Cari user di database
	err := userCollection.FindOne(ctx, bson.M{"username": input.Username}).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Cek Password
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Buat JWT Token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.MapClaims{
		"username": foundUser.Username,
		"exp":      expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Ambil Secret Key dari Env (atau default)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "rahasia_super_aman" // Ganti ini nanti di production!
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}