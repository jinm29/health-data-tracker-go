package services


import (
	"net/http" 

	"github.com/Fidel-wole/wearable-integration/db"
	sqlc "github.com/Fidel-wole/wearable-integration/db/sqlc"
	"github.com/gin-gonic/gin"
)

func InsertWearableDevice(c *gin.Context){
	var data sqlc.AddWearableDevicesParams
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
    queries := db.GetQueries()

	device, err := queries.AddWearableDevices(c, data)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device added successfully", "data":device})
}

func GetAvailableDevices(c *gin.Context){
	queries := db.GetQueries()
	devices, err := queries.GetAvailableDevices(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Devices fetched successfully", "data":devices})
}

func AssignDeviceToPatient(c *gin.Context) {
	var data sqlc.AssignDeviceToPatientParams
	var assignTo sqlc.AssignDeviceParams

	// Bind JSON input to the data structure
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queries := db.GetQueries()

	// Assign device to patient
	if err := queries.AssignDeviceToPatient(c, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	assignTo = sqlc.AssignDeviceParams{
		AssignedPatientID: &data.ID, 
		ID:                data.AssignedDeviceID,
	}

	
	if err := queries.AssignDevice(c, assignTo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Device assigned to patient successfully"})
}
