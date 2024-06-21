package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/PoseidonPSE2/code_backend/database"
	"github.com/gin-gonic/gin"
)

// @Summary Show all users
// @Description Get all users
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} database.User
// @Router /users [get]
func GetUsers(c *gin.Context) {
	log.Print(db)
	idStr := c.Query("id")
	if idStr == "" {
		var users []database.User
		result := db.Find(&users)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var user database.User
		result := db.First(&user, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

// @Summary Create a user
// @Description Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body database.User true "User"
// @Success 201 {object} database.User
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var user database.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// @Summary Update a user
// @Description Update an existing user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body database.User true "User"
// @Success 200 {object} database.User
// @Router /users [put]
func UpdateUser(c *gin.Context) {
	var user database.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary Delete a user
// @Description Delete an existing user
// @Tags Users
// @Accept json
// @Produce json
// @Param id query int true "User ID"
// @Success 204
// @Router /users [delete]
func DeleteUser(c *gin.Context) {
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
	result := db.Delete(&database.User{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
