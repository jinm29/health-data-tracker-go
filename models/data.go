package models

type Data struct {
    PatientID   string `json:"patient_id"`
    DeviceID    string  `json:"device_id"`
    HeartRate   float64 `json:"heart_rate"`
    OxygenLevel float64 `json:"oxygen_level"`
}
