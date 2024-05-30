package api

import (
	"math"
	"net/http"
	"strconv"

	"github.com/PoseidonPSE2/code_backend/database"
	"github.com/gin-gonic/gin"
)

// RefillStationResponse represents a refill station in the response
type RefillStationResponse struct {
	ID                uint    `json:"id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Address           string  `json:"address"`
	WaterSource       string  `json:"water_source"`
	OpeningTimes      string  `json:"opening_times"`
	Active            bool    `json:"active"`
	Type              string  `json:"type"`
	OfferedWaterTypes string  `json:"offered_water_types"`
	ImagePath         *string `json:"image_path"`
}

// RefillStationMarkerResponse represents a refill station marker in the response
type RefillStationMarkerResponse struct {
	ID        uint    `json:"id"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Status    bool    `json:"status"`
}

// CreateRefillStationRequest represents a request to create a refill station
type CreateRefillStationRequest struct {
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Address           string  `json:"address"`
	WaterSource       string  `json:"water_source"`
	OpeningTimes      string  `json:"opening_times"`
	Active            bool    `json:"active"`
	Type              string  `json:"type"`
	OfferedWaterTypes string  `json:"offered_water_types"`
	ImagePath         *string `json:"image_path"`
}

// UpdateRefillStationRequest represents a request to update a refill station
type UpdateRefillStationRequest struct {
	ID                uint    `json:"id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Address           string  `json:"address"`
	WaterSource       string  `json:"water_source"`
	OpeningTimes      string  `json:"opening_times"`
	Active            bool    `json:"active"`
	Type              string  `json:"type"`
	OfferedWaterTypes string  `json:"offered_water_types"`
	ImagePath         *string `json:"image_path"`
}

// @Summary Show all refill stations
// @Description Get all refill stations
// @Tags refill_stations
// @Accept  json
// @Produce  json
// @Success 200 {array} RefillStationResponse
// @Router /refill_stations [get]
func GetRefillStations(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		var stations []database.RefillStation
		result := db.Find(&stations)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, stations)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var station database.RefillStation
		result := db.First(&station, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, station)
	}
}

// @Summary Get all refill station markers
// @Description Get all refill station markers with specific attributes
// @Tags refill_stations
// @Accept  json
// @Produce  json
// @Success 200 {array} RefillStationMarkerResponse
// @Router /refill_stations/markers [get]
func GetAllRefillstationMarker(c *gin.Context) {
	var stations []database.RefillStation
	result := db.Select("id, longitude, latitude, active").Find(&stations)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var markers []RefillStationMarkerResponse
	for _, station := range stations {
		markers = append(markers, RefillStationMarkerResponse{
			ID:        station.ID,
			Longitude: station.Longitude,
			Latitude:  station.Latitude,
			Status:    station.Active,
		})
	}

	respondWithJSON(c, http.StatusOK, markers)
}

// @Summary Get a refill station by ID
// @Description Get a refill station by its ID
// @Tags refill_stations
// @Accept  json
// @Produce  json
// @Param id path int true "Refill Station ID"
// @Success 200 {object} RefillStationResponse
// @Router /refill_stations/{id} [get]
func GetRefillStationById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var station database.RefillStation
	result := db.First(&station, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusOK, station)
}

// @Summary Get the average review score for a refill station
// @Description Get the average review score for a refill station by its ID
// @Tags refill_station_reviews
// @Accept  json
// @Produce  json
// @Param id query int true "Refill Station ID"
// @Success 200 {number} float64
// @Router /refill_stations/{id}/reviews [get]
func GetRefillStationReviewByID(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var reviews []database.RefillStationReview
	result := db.Where("station_id = ?", id).Find(&reviews)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if len(reviews) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No reviews found"})
		return
	}

	var totalCleanness, totalAccessibility, totalWaterQuality float64
	for _, review := range reviews {
		totalCleanness += float64(review.Cleanness)
		totalAccessibility += float64(review.Accessibility)
		totalWaterQuality += float64(review.WaterQuality)
	}

	average := (totalCleanness + totalAccessibility + totalWaterQuality) / (float64(len(reviews)) * 3)
	average = math.Round(average*10) / 10 // Round to 1 decimal place

	c.JSON(http.StatusOK, gin.H{"average": average})
}

// @Summary Create a refill station
// @Description Create a new refill station
// @Tags refill_stations
// @Accept  json
// @Produce  json
// @Param station body CreateRefillStationRequest true "Refill Station"
// @Success 201 {object} RefillStationResponse
// @Router /refill_stations [post]
func CreateRefillStation(c *gin.Context) {
	var station database.RefillStation
	if err := c.ShouldBindJSON(&station); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Create(&station)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusCreated, station)
}

// @Summary Update a refill station
// @Description Update an existing refill station
// @Tags refill_stations
// @Accept  json
// @Produce  json
// @Param station body UpdateRefillStationRequest true "Refill Station"
// @Success 200 {object} RefillStationResponse
// @Router /refill_stations [put]
func UpdateRefillStation(c *gin.Context) {
	var station database.RefillStation
	if err := c.ShouldBindJSON(&station); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Save(&station)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusOK, station)
}

// @Summary Delete a refill station
// @Description Delete an existing refill station
// @Tags refill_stations
// @Accept  json
// @Produce  json
// @Param id query int true "Refill Station ID"
// @Success 204
// @Router /refill_stations [delete]
func DeleteRefillStation(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	result := db.Delete(&database.RefillStation{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
