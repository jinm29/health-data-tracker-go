package main

import (
	"log" // Import log for error logging

	"github.com/Fidel-wole/wearable-integration/config"
	"github.com/Fidel-wole/wearable-integration/db"
	"github.com/Fidel-wole/wearable-integration/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database and handle errors
	db.InitDB()

	r := gin.Default()


	config.InitInfluxDB()


	go config.SetupMQTT()


	r.POST("/add-patient", services.AddPatient)
	r.GET("/patients", services.GetAllPatients)
	r.GET("/patient/:user_id", services.GetSinglePatientData)
	r.POST("/add-device", services.InsertWearableDevice)
	r.GET("/available-devices", services.GetAvailableDevices)
	r.PUT("/assign-device", services.AssignDeviceToPatient)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Could not run the server: %v", err)
	}
}
