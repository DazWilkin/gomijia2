package main

import (
	"flag"
	"log"

	"github.com/currantlabs/ble/linux"
)

var (
	configFile = flag.String("config_file", "config.ini", "Config file location")
)

func main() {
	flag.Parse()

	log.Print("[main] Reading configuration")
	config, err := NewConfig(*configFile)
	if err != nil {
		log.Fatal("Unable to parse configuration")
	}

	log.Print("[main] Starting Linux Device")
	config.Host, err = linux.NewDevice()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[main] MQTT broker: %s", config.MQTT.Server())
	if err := config.MQTT.Connect("xiaomi"); err != nil {
		log.Print("[main] Unable to connect to MQTT broker")
	}

	for _, device := range config.Devices {

		log.Printf("[main:%s] Dialing (%s)", device.Name, device.Addr)
		if err := device.Connect(config.Host); err != nil {
			log.Printf("[main:%s] Failed to connect to device", device.Name)
			continue
		}

		log.Printf("[main:%s] Registering handler", device.Name)
		device.RegisterHandler(config.MQTT)

	}

	// Loop forever while notification handler respond
	for {
	}

	// for _, device := range config.Devices {
	// 	if err := device.Disconnect(); err != nil {
	// 		log.Printf("[main:%s] Failed to disconnect from device", device.Name)
	// 		continue
	// 	}

	// }

	// log.Print("[main] Stopping Linux Device")
	// config.Host.Stop()

}
