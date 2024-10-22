package controller

import (
    "net/http"

    "github.com/Fidel-wole/wearable-integration/services"
    "github.com/gin-gonic/gin"
)

// HealthDataController handles health data-related requests
type HealthDataController struct {
    DataService *services.DataService
}

// NewHealthDataController initializes a new HealthDataController
func NewHealthDataController(ds *services.DataService) *HealthDataController {
    return &HealthDataController{
        DataService: ds,
    }
}

// GetHealthDataByDeviceID handles GET requests to fetch health data by device ID
func (hd *HealthDataController) GetHealthDataByDeviceID(c *gin.Context) {
    patientId := c.Param("patient_id")

    // Set context for the query
    ctx := c.Request.Context()

    // Fetch health data
    data, err := hd.DataService.GetHealthDataByDeviceID(ctx, patientId)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Check if no data is found
    if len(data) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"message": "No data found for the specified device ID"})
        return
    }

    // Respond with the health data
    c.JSON(http.StatusOK, data)
}
