# Health Data Tracker

A Go-based application designed to integrate wearable device data using MQTT and InfluxDB. This project allows for real-time collection and storage of health-related metrics such as heart rate and oxygen levels.

## Features
1. Real-time Data Processing: Collects health data from wearable devices via MQTT.
2. InfluxDB Integration: Saves processed data into InfluxDB for efficient time-series data storage and analysis.
3. Graceful Shutdown: Handles system signals for graceful application termination.

## Tech Stack
Go
MQTT (Eclipse Paho)
InfluxDB

## Installation
Prerequisites:

Go (version 1.16 or higher)

InfluxDB

MQTT Broker (e.g., Mosquitto)

## Clone the repository

```bash
git clone https://github.com/Fidel-wole/health_data_tracker.git
cd health_data_tracker
```

## Set up environment variables
Create a .env file in the root directory and add your InfluxDB token:
```plaintext
INFLUXDB_TOKEN=your_influxdb_token
```
## Install Go dependencies:
```bash
go run main.go
```
Run the MQTT broker (if you haven't done so already):

If using Mosquitto, you can install it via Homebrew (on macOS) or download it from Mosquitto's official website

Run the application:
```bash
go run main.go
```
## Usage
The application listens for MQTT messages published to the wearable/data topic.

Data format (JSON) expected in the payload:
```json
{
    "device_id": "device_199",
    "heart_rate": 175,
    "oxygen_level": 75
}
```
The application processes the incoming data and saves it to InfluxDB.
## Testing
To test the MQTT integration, you can publish test messages using an MQTT client or command-line tool. For example, using mosquitto_pub:
```bash
mosquitto_pub -h localhost -t wearable/data -m '{"device_id": "device_199", "heart_rate": 175, "oxygen_level": 75}'
```
## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
