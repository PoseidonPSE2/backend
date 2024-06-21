package api

import (
	"encoding/base64"
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

// DecodeBase64ToBytes takes a base64 encoded string and returns a byte array
func DecodeBase64ToBytes(encodedStr string) ([]byte, error) {
	// Decode the base64 string
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedStr)
	if err != nil {
		return nil, err
	}
	return decodedBytes, nil
}
