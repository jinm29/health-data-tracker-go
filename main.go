package main

import (
    "github.com/Fidel-wole/wearable-integration/config"

)

func main() {

    config.InitInfluxDB()
    go config.SetupMQTT()

    select {}
}
