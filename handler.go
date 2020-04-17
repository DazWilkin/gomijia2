package main

import (
	"encoding/hex"
	"log"
)

func handlerPublisher(mqtt MQTT, name string) func(req []byte) {
	return func(req []byte) {
		s := hex.EncodeToString(req)
		r, err := Unmarshall(req)
		if err != nil {
			log.Printf("[handler:%s] Unable to unmarshal data (%s)", name, s)
		}
		log.Printf("[handler:%s] %s (%s)", name, r.String(), s)
		mqtt.Publish(name, r)
	}
}
