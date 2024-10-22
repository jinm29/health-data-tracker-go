package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/Fidel-wole/wearable-integration/config"
	"github.com/Fidel-wole/wearable-integration/controller"
	"github.com/Fidel-wole/wearable-integration/db"
	"github.com/Fidel-wole/wearable-integration/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow requests from any origin
    },
}

var (
    userConnections = make(map[string]map[string]*websocket.Conn)
    mu              sync.Mutex
)

func handleWebSocket(c *gin.Context) {
    userID := c.Param("user_id")
    connectionID := c.Param("connection_id")

    // Check if the user is already connected
    mu.Lock()
    if _, ok := userConnections[userID]; !ok {
        userConnections[userID] = make(map[string]*websocket.Conn)
    }
    mu.Unlock()

    // Establish the WebSocket connection
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("Failed to upgrade connection for user %s: %v", userID, err)
        return
    }
    defer conn.Close()

    mu.Lock()
    userConnections[userID][connectionID] = conn
    log.Printf("User  %s connected via WebSocket. Total connections: %d", userID, len(userConnections[userID]))
    logConnectedUsers()
    mu.Unlock()

    // Listen for messages
    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Printf("WebSocket read error for user %s: %v", userID, err)
            break // Exit the loop on error
        }
        log.Printf("Received message from user %s: %s", userID, message)
    }

    // Handle cleanup after connection closure
    mu.Lock()
    delete(userConnections[userID], connectionID)
    log.Printf("Connection closed for user %s. Remaining connections: %d", userID, len(userConnections[userID]))
    logConnectedUsers()
    mu.Unlock()
}

func logConnectedUsers() {
    mu.Lock()
    defer mu.Unlock()

    log.Println("Currently connected users:")
    for userID := range userConnections {
        log.Println(userID)
    }
}

func main() {
    godotenv.Load()
    // Initialize the database
    db.InitDB()

    r := gin.Default()

    // CORS middleware configuration
    r.Use(cors.New(cors.Config{
        AllowAllOrigins:  true, // Allows all origins
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // Allowed HTTP methods
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allowed headers
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true, // Allow credentials (cookies, authorization headers)
    }))

    influxDBClient := config.InitInfluxDB()

    go config.SetupMQTT(userConnections)
    // Initialize InfluxDB


    // Initialize the DataService and HealthDataController
    dataService := &services.DataService{InfluxDBClient: influxDBClient}
    healthDataController := controller.NewHealthDataController(dataService)
    // Routes
    r.POST("/add-patient", services.AddPatient)
    r.GET("/patients", services.GetAllPatients)
    r.GET("/patient/:user_id", services.GetSinglePatientData)
    r.POST("/add-device", services.InsertWearableDevice)
    r.GET("/available-devices", services.GetAvailableDevices)
    r.PUT("/assign-device", services.AssignDeviceToPatient)
    r.GET("/health-stats/:patient_id", healthDataController.GetHealthDataByDeviceID)
    // WebSocket route
    r.GET("/ws/:user_id", handleWebSocket)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Default to 8080 if not set
    }

    // Start the server
    if err := r.Run(":" + port); err != nil {
        log.Fatalf("Could not run the server: %v", err)
    }
}
