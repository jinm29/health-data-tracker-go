// services/data_service.go
package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Fidel-wole/wearable-integration/models"
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
)


type DataService struct {
    InfluxDBClient *influxdb3.Client
}

// SaveDataToInfluxDB remains unchanged
func (ds *DataService) SaveDataToInfluxDB(data models.Data) {
    // Validate the data before proceeding
    if data.HeartRate < 0 || data.OxygenLevel < 0 {
        log.Printf("Invalid data for device %s: HeartRate=%d, OxygenLevel=%d", data.DeviceID, data.HeartRate, data.OxygenLevel)
        return
    }

    // Create the point
    point := influxdb3.NewPoint("health_data",
        map[string]string{
            "device_id": data.DeviceID,
            "patient_id": data.PatientID,
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

    if err := ds.InfluxDBClient.WritePoints(ctx, points, influxdb3.WithDatabase("wearable")); err != nil {
        log.Printf("Error writing to InfluxDB for device %s: %v", data.DeviceID, err)
    } else {
        log.Println("Data saved to InfluxDB")
    }
}

// New method to fetch health data from InfluxDB
func (ds *DataService) GetHealthDataByDeviceID(ctx context.Context, patientId string) ([]models.Data, error) {
    query := fmt.Sprintf(`SELECT heart_rate, oxygen_level 
                          FROM "health_data" 
                          WHERE patient_id = '%s'`, patientId)

    // Set query options
    queryOptions := influxdb3.QueryOptions{
        Database: "wearable", // Set the database you want to query
    }

    // Execute the query
    iterator, err := ds.InfluxDBClient.QueryWithOptions(ctx, &queryOptions, query)
    if err != nil {
        return nil, fmt.Errorf("query error: %v", err)
    }
    
    // Slice to hold the results
    var results []models.Data

    // Iterate through the results
    for iterator.Next() {
        // Get the current record's value
        value := iterator.Value()
        
        // Access the data fields
        heartRate, ok1 := value["heart_rate"].(float64)
        oxygenLevel, ok2 := value["oxygen_level"].(float64)
        deviceId, ok3 := value["device_id"].(string)
        // Check if the data is valid
        if ok1 && ok2 && ok3{
            data := models.Data{
                HeartRate:   heartRate,
                OxygenLevel: oxygenLevel,
                DeviceID: deviceId,
                PatientID:    patientId, 
            }
            results = append(results, data)
        }
    }

    return results, nil
}

