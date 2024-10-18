package config

import (
    "github.com/InfluxCommunity/influxdb3-go/influxdb3"
    mqtt "github.com/eclipse/paho.mqtt.golang"
    "log"
    "os"
)

var MQTTClient mqtt.Client
var InfluxDBClient *influxdb3.Client

func InitInfluxDB() {
    
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
    log.Println("Connected to InfluxDB")
}

func SetupMQTT() {
    opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
    opts.SetClientID("wearable-client")
    MQTTClient = mqtt.NewClient(opts)
    if token := MQTTClient.Connect(); token.Wait() && token.Error() != nil {
        log.Fatal(token.Error())
    }
    log.Println("MQTT client connected")
}
