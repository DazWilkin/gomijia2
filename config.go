package main

import (
	"fmt"
	"log"

	"github.com/currantlabs/ble/linux"
	"gopkg.in/ini.v1"
)

// Config represents a configuration
type Config struct {
	MQTT    MQTT
	Devices []Device
	Host    *linux.Device
}

// NewConfig returns a new Config
func NewConfig(file string) (*Config, error) {
	log.Printf("[Config] Loading Configuration (%s)", file)
	cfg, err := ini.Load(file)
	if err != nil {
		return &Config{}, err
	}

	sec, err := cfg.GetSection("MQTT")
	if err != nil {
		return &Config{}, err
	}
	if !sec.HasKey("host") {
		return &Config{}, fmt.Errorf("Configuration requires MQTT host")
	}

	if !sec.HasKey("port") {
		log.Print("MQTT port not defined; defaulting to 1883")
		sec.NewKey("port", "1883")
	}

	mqtt := MQTT{
		Host: sec.Key("host").String(),
		Port: sec.Key("port").String(),
		User: sec.Key("user").String(),
		Pass: sec.Key("pass").String(),
	}

	sec, err = cfg.GetSection("Devices")
	if err != nil {
		return &Config{}, err
	}
	names := sec.KeyStrings()

	devices := []Device{}
	for i, name := range names {
		addr := sec.Key(name).String()
		log.Printf("[Config] Device %02d: %s (%s)", i, name, addr)
		devices = append(devices, Device{
			Name: name,
			Addr: addr,
		})
	}

	return &Config{
		MQTT:    mqtt,
		Devices: devices,
	}, nil
}
