// models/place_model.go
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// GeoJSON mendefinisikan struktur untuk data geografis Point.
type GeoJSON struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"` // Format: [longitude, latitude]
}

// Place adalah model utama untuk data lokasi kita.
type Place struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Location GeoJSON            `json:"location" bson:"location"`
}