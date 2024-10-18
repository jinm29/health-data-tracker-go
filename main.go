package main

import (
    "github.com/Fidel-wole/wearable-integration/config"
    "github.com/Fidel-wole/wearable-integration/controllers"
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    // Initialize configurations
    config.InitInfluxDB()
    config.SetupMQTT()

    
    router.POST("/ingest", controllers.IngestData)

    // Start the server
    router.Run(":8080")
}
