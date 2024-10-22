// config/config.go
package config

import (
	"encoding/json"
	//"fmt"
	"log"
	"os"

	"github.com/Fidel-wole/wearable-integration/models"
	"github.com/Fidel-wole/wearable-integration/services"
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
	// "github.com/joho/godotenv"
)

var MQTTClient mqtt.Client
var InfluxDBClient *influxdb3.Client
var DataService *services.DataService

func InitInfluxDB() {
   // godotenv.Load()
    url := "https://us-east-1-1.aws.cloud2.influxdata.com"
    token := os.Getenv("INFLUXDB_TOKEN")

    client, err := influxdb3.New(influxdb3.ClientConfig{
        Host:  url,
        Token: token,
    })
    if err != nil {
        log.Fatalf("Failed to create InfluxDB client: %v", err)
    }

    InfluxDBClient = client
    DataService = &services.DataService{InfluxDBClient: InfluxDBClient}
    log.Println("Connected to InfluxDB")
}

func SetupMQTT(userConnections map[string]map[string]*websocket.Conn) {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("wearable-client")

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		data := models.Data{}
		log.Printf("Raw payload: %s", string(msg.Payload()))

		// Unmarshal the JSON payload into `data`
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			log.Printf("Failed to unmarshal MQTT message: %v", err)
			return
		}

		// Save the data to InfluxDB
		DataService.SaveDataToInfluxDB(data)
		log.Println("Data received from MQTT and stored.")

		// Check if there are WebSocket connections for the specific user
		if connections, ok := userConnections[data.PatientID]; ok {
			for _, conn := range connections {
				// Send data to each user's WebSocket connection
				if err := conn.WriteJSON(data); err != nil {
					log.Printf("Failed to send data to user %s: %v", data.PatientID, err)
					conn.Close()

				} else {
					log.Printf("Data sent to user %s via WebSocket", data.PatientID)
				}
			}
		} else {
			log.Printf("No active WebSocket connection for user %s", data.PatientID)
		}
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	log.Println("MQTT client connected and listening for messages")

	// Subscribe to the topic where the device sends data
	if token := client.Subscribe("wearable/data", 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	log.Println("Subscribed to MQTT topic 'wearable/data'")
}
