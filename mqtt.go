package main

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTT represents an MQTT client configuration
type MQTT struct {
	Host   string
	Port   string
	User   string
	Pass   string
	Client mqtt.Client
	id     string
}

// Server returns a new MQTT connection string
func (m *MQTT) Server() string {
	return fmt.Sprintf("tcp://%s:%s", m.Host, m.Port)
}

// Connect connects to an MQTT broker
func (m *MQTT) Connect(id string) error {
	log.Printf("[MQTT:Connect] Connecting: %s", id)

	m.id = id

	opts := mqtt.NewClientOptions().AddBroker(m.Server())
	opts.SetClientID(id)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	if m.User != "" {
		opts.SetUsername(m.User)
	}
	if m.Pass != "" {
		opts.SetPassword(m.Pass)
	}

	log.Print("[MQTT] Creating")
	m.Client = mqtt.NewClient(opts)

	log.Print("[MQTT] Connecting")
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		return (token.Error())
	}
	return nil
}

// Disconnect disconnects from an MQTT broker
func (m *MQTT) Disconnect() {
	log.Print("[MQTT] Disconnecting")
	m.Client.Disconnect(250)
}

// Publish publishes a message to an MQTT topic
func (m *MQTT) Publish(name string, r *Reading) {
	log.Printf("[MQTT] Publishing %s (%s)", name, r.String())
	format := "prometheus/job/%s/node/%s/%s"

	// Publish Temperature
	m.Client.Publish(fmt.Sprintf(format, m.id, name, "temperature"), 0, false, fmt.Sprintf("%04f", r.Temperature))

	// Publish Humidity
	m.Client.Publish(fmt.Sprintf(format, m.id, name, "humidity"), 0, false, fmt.Sprintf("%04f", r.Humidity))

}
