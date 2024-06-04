package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func respondWithJSON(c *gin.Context, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(status, "application/json", response)
}

func SetDB(new_db *gorm.DB) {
	db = new_db
}
