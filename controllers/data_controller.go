package controllers

import (
    "net/http"
    "github.com/Fidel-wole/wearable-integration/models"
    "github.com/Fidel-wole/wearable-integration/services"
    "github.com/gin-gonic/gin"
)

func IngestData(c *gin.Context) {
    var data models.Data
    if err := c.ShouldBindJSON(&data); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    services.ProcessAndStoreData(data)
    c.JSON(http.StatusOK, gin.H{"message": "Data received"})
}
