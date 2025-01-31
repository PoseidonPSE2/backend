package api

import (
	"net/http"
	"strconv"

	"github.com/PoseidonPSE2/code_backend/database"
	"github.com/gin-gonic/gin"
)

type StationReviewAverage struct {
	Cleanness    float64 `json:"cleanness"`
	Accesibility float64 `json:"accesibility"`
	WaterQuality float64 `json:"waterQuality"`
}

type StationImage struct {
	StationImage []byte `json:"station_image"`
}

// @Summary Show all refill stations
// @Description Get all refill stations
// @Tags Refill Stations
// @Accept json
// @Produce json
// @Success 200 {array} database.RefillStation
// @Router /refill_stations [get]
func GetRefillStations(c *gin.Context) {
	var stations []database.RefillStation
	result := db.Find(&stations)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, stations)
}

// @Summary Get all refill station markers
// @Description Get all refill station markers with specific attributes
// @Tags Refill Stations
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /refill_stations/markers [get]
func GetAllRefillstationMarker(c *gin.Context) {
	var stations []database.RefillStation
	result := db.Select("id, longitude, latitude, active").Find(&stations)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var markers []map[string]interface{}
	for _, station := range stations {
		markers = append(markers, map[string]interface{}{
			"id":        station.ID,
			"longitude": station.Longitude,
			"latitude":  station.Latitude,
			"status":    station.Active,
		})
	}

	c.JSON(http.StatusOK, markers)
}

// @Summary Get a refill station by ID
// @Description Get a refill station by its ID
// @Tags Refill Stations
// @Accept json
// @Produce json
// @Param id path int true "Refill Station ID"
// @Success 200 {object} database.RefillStation
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
	c.JSON(http.StatusOK, station)
}

// @Summary Get the image from a refill station by ID
// @Description Get the image from a refill station by ID
// @Tags Refill Stations
// @Accept json
// @Produce json
// @Param id path int true "Refill Station ID"
// @Success 200 {object} StationImage
// @Router /refill_stations/image/{id} [get]
func GetRefillStationImageById(c *gin.Context) {
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

	if station.RefillStationImage == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refill Station has no image"})
		return
	}

	byteArray, err := DecodeBase64ToBytes(*station.RefillStationImage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding base64 string"})
		return
	}

	response := StationImage{
		StationImage: byteArray,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get the average review score for a refill station
// @Description Get the average review score for a refill station by its ID
// @Tags Refill Stations
// @Accept json
// @Produce json
// @Param id path int true "Refill Station ID"
// @Success 200 {object} StationReviewAverage
// @Router /refill_stations/{id}/reviews [get]
func GetRefillStationReviewsAverageByID(c *gin.Context) {
	var station database.RefillStation
	var stationReviews []database.RefillStationReview
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	stationQueryResult := db.First(&station, id)
	if stationQueryResult.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": stationQueryResult.Error.Error()})
		return
	}

	reviewsQueryResult := db.Where("station_id = ?", id).Find(&stationReviews)
	if reviewsQueryResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": reviewsQueryResult.Error.Error()})
		return
	}

	if len(stationReviews) == 0 {
		c.JSON(http.StatusOK, gin.H{"cleanness": 0.0, "accesibility": 0.0, "waterQuality": 0.0})
		return
	}

	var totalCleanness, totalAccessibility, totalWaterQuality float64
	for _, review := range stationReviews {
		totalCleanness += float64(review.Cleanness)
		totalAccessibility += float64(review.Accessibility)
		totalWaterQuality += float64(review.WaterQuality)
	}

	amountReviews := (len(stationReviews))

	cleannessAverage := totalCleanness / float64(amountReviews)
	accessibilityAverage := totalAccessibility / float64(amountReviews)
	waterQualityAverage := totalWaterQuality / float64(amountReviews)

	response := StationReviewAverage{
		Cleanness:    cleannessAverage,
		Accesibility: accessibilityAverage,
		WaterQuality: waterQualityAverage,
	}

	respondWithJSON(c, http.StatusOK, response)
}

// @Summary Create a refill station
// @Description Create a new refill station
// @Tags Refill Stations
// @Accept json
// @Produce json
// @Param station body database.RefillStation true "Refill Station"
// @Success 201 {object} database.RefillStation
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
	c.JSON(http.StatusCreated, station)
}

// @Summary Update a refill station
// @Description Update an existing refill station
// @Tags Refill Stations
// @Accept  json
// @Produce  json
// @Param station body database.RefillStation true "Refill Station"
// @Success 200 {object} database.RefillStation
// @Router /refill_stations [put]
func UpdateRefillStation(c *gin.Context) {
	var requestStation database.RefillStation
	if err := c.ShouldBindJSON(&requestStation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var helpStation database.RefillStation
	if result := db.First(&helpStation, requestStation.ID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Refill Station with ID not found"})
		return
	}
	db.Model(&helpStation).Updates(requestStation)
	c.JSON(http.StatusOK, requestStation)
}

// @Summary Delete a refill station
// @Description Delete an existing refill station
// @Tags Refill Stations
// @Accept json
// @Produce json
// @Param id path int true "Refill Station ID"
// @Success 204
// @Router /refill_stations/{id} [delete]
func DeleteRefillStation(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var station database.RefillStation
	if result := db.First(&station, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Refill Station with ID not found"})
		return
	}

	result := db.Delete(&database.RefillStation{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
