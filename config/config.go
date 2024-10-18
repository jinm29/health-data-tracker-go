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
	"github.com/joho/godotenv"
)

var MQTTClient mqtt.Client
var InfluxDBClient *influxdb3.Client
var DataService *services.DataService

func InitInfluxDB() {
    godotenv.Load()
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

func SetupMQTT() {
    opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
    opts.SetClientID("wearable-client")

    opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
        // Process the incoming MQTT message
        data := models.Data{}
        log.Printf("Raw payload: %s", string(msg.Payload()))
        // Unmarshal the JSON payload into `data`
        if err := json.Unmarshal(msg.Payload(), &data); err != nil {
            log.Printf("Failed to unmarshal MQTT message: %v", err)
            return
        }
        log.Printf("Data after unmarshal: %+v", data)
        // Process and store the data in InfluxDB
        DataService.SaveDataToInfluxDB(data)
        log.Println("Data received from MQTT and processed.")
        log.Printf("The processed data is: %+v", data)
    })

    MQTTClient = mqtt.NewClient(opts)
    if token := MQTTClient.Connect(); token.Wait() && token.Error() != nil {
        log.Fatal(token.Error())
    }
    log.Println("MQTT client connected and listening for messages")

    // Subscribe to the topic where the device sends data
    if token := MQTTClient.Subscribe("wearable/data", 0, nil); token.Wait() && token.Error() != nil {
        log.Fatal(token.Error())
    }
    log.Println("Subscribed to MQTT topic 'wearable/data'")
}
