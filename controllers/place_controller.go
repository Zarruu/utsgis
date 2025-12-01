// controllers/place_controller.go
package controllers

import (
	"context"
	"net/http"
	"time"

	// "go-mongo-geojson/config"
	"go-mongo-geojson/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
)

var placeCollection *mongo.Collection

// FUNGSI BARU: untuk menginisialisasi collection dari main.go
func InitPlaceCollection(client *mongo.Client) {
	// Ganti "your_db_name" dengan nama database Anda
	placeCollection = client.Database("location_db").Collection("places")
}

// CreatePlace: Menambah data tempat baru
func CreatePlace(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var place models.Place
	if err := c.ShouldBindJSON(&place); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pastikan tipe GeoJSON adalah "Point"
	place.Location.Type = "Point"

	result, err := placeCollection.InsertOne(ctx, place)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create place"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}

// GetPlaces: Mendapatkan semua data tempat
func GetPlaces(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var places []models.Place
	cursor, err := placeCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get places"})
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &places); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, places)
}

// UpdatePlace: Memperbarui data tempat
func UpdatePlace(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	placeId := c.Param("id")
	objId, _ := primitive.ObjectIDFromHex(placeId)

	var place models.Place
	if err := c.ShouldBindJSON(&place); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{"$set": bson.M{
		"name":     place.Name,
		"location": place.Location,
	}}

	result, err := placeCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update place"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Place not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Place updated successfully"})
}

// DeletePlace: Menghapus data tempat
func DeletePlace(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	placeId := c.Param("id")
	objId, _ := primitive.ObjectIDFromHex(placeId)

	result, err := placeCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete place"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Place not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Place deleted successfully"})
}