package main

import (
	"encoding/binary"
	"fmt"
	"log"
)

// Reading represents a Temperature|Humidity readings
type Reading struct {
	Temperature float64
	Humidity    float64
    Battery     float64
}

// Equal determine whether one Reading is equal to another
func (r *Reading) Equal(s Reading) bool {
    log.Printf("[Reading:Equal] Temperatures: %f (%f); Humidity: %f (%f); Battery: %f (%f)", r.Temperature, s.Temperature, r.Humidity, s.Humidity, r.Battery, s.Battery)
	if r.Temperature != s.Temperature {
		log.Print("[Reading:Equal] Temperatures don't match")
	}
	if r.Humidity != s.Humidity {
		log.Print("[Reading:Equal] Humidities don't match")
	}
	if r.Battery != s.Battery {
		log.Print("[Reading:Equal] Battery voltages don't match")
	}
	return r.Temperature == s.Temperature && r.Humidity == s.Humidity && r.Battery == s.Battery
}

// ToString converts a Reading to a string
func (r *Reading) String() string {
	return fmt.Sprintf("Temperature: %.04f; Humidity: %.04f; Battery: %.04f", r.Temperature, r.Humidity, r.Battery)
}

// NewReading returns a new Reading
func NewReading(t, h, v float64) *Reading {
	return &Reading{
		Temperature: t,
		Humidity:    h,
		Battery:    v,
	}
}

// Unmarshall converts an encoded reading into a Reading
func Unmarshall(req []byte) (*Reading, error) {
	// 00 01 02 03 04
	// T2 T1 HX V2 V1
	l := len(req)
	if l != 5 {
		log.Printf("[X] Expecting 5 bytes; got %d", l)
		return &Reading{}, fmt.Errorf("Expecting 5 bytes got %d", l)
	}
	// Temperature is stored little endian *100
	t := float64(int(binary.LittleEndian.Uint16(req[0:2]))) / 100.0
	h := float64(req[2]) / 100.0
	// Battery voltage is stored little endian * 1000
	v := float64(int(binary.LittleEndian.Uint16(req[3:5]))) / 1000.0
	return &Reading{
		Temperature: t,
		Humidity:    h,
        Battery:     v,
	}, nil
}
