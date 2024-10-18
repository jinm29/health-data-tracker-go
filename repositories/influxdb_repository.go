package repositories

import (
    "context"
    "log"
    "time"
    "github.com/Fidel-wole/wearable-integration/config"
    "github.com/Fidel-wole/wearable-integration/models"
    "github.com/InfluxCommunity/influxdb3-go/influxdb3"
)

func SaveDataToInfluxDB(data models.Data) {
    // Validate the data before proceeding
    if data.HeartRate < 0 || data.OxygenLevel < 0 {
        log.Printf("Invalid data for device %s: HeartRate=%d, OxygenLevel=%d", data.DeviceID, data.HeartRate, data.OxygenLevel)
        return
    }

    // Create the point
    point := influxdb3.NewPoint("health_data",
        map[string]string{
            "device_id": data.DeviceID,
        },
        map[string]interface{}{
            "heart_rate":   data.HeartRate,
            "oxygen_level": data.OxygenLevel,
        },
        time.Now(),
    )

    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Create a slice of points
    points := []*influxdb3.Point{point}

    if err := config.InfluxDBClient.WritePoints(ctx, points, influxdb3.WithDatabase("wearable")); err != nil {
        log.Printf("Error writing to InfluxDB for device %s: %v", data.DeviceID, err)
    } else {
        log.Println("Data saved to InfluxDB")
    }
}

