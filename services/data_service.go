package services

import (
    "github.com/Fidel-wole/wearable-integration/models"
    "github.com/Fidel-wole/wearable-integration/repositories"
    "fmt"
)

func ProcessAndStoreData(data models.Data) {
    checkForAlerts(data)
    repositories.SaveDataToInfluxDB(data)
}

func checkForAlerts(data models.Data) {
    if data.HeartRate > 100 {
        fmt.Println("Alert: High heart rate detected!")
    }
    if data.OxygenLevel < 90 {
        fmt.Println("Alert: Low oxygen level detected!")
    }
}
