package services

import (
	"net/http" // Added for HTTP status codes
	"strconv"

	"github.com/Fidel-wole/wearable-integration/db"
	sqlc "github.com/Fidel-wole/wearable-integration/db/sqlc"
	"github.com/gin-gonic/gin"
)


func AddPatient(c *gin.Context) {
	var data sqlc.CreatePatientParams

	// Bind JSON input to the data structure
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the database queries
	queries := db.GetQueries()

	// Create a new patient
	createdPatient, err := queries.CreatePatient(c, data); 
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Patient added successfully", "data":createdPatient})
}

func GetAllPatients (c *gin.Context){
	queries := db.GetQueries()

	result, err :=queries.GetAllPatients(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Patients fetcched successfully", "data":result})
}

func GetSinglePatientData(c *gin.Context) {
    // Extract user ID from URL parameters
    userID := c.Param("user_id")
    

    userIDInt, err := strconv.Atoi(userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

	userIDInt32 := int32(userIDInt)
    
    queries := db.GetQueries()
    
    result, err := queries.GetPatientWithDevice(c, userIDInt32)
    
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Return the result
    c.JSON(http.StatusOK, result)
}
