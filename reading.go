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
}

// Equal determine whether one Reading is equal to another
func (r *Reading) Equal(s Reading) bool {
	log.Printf("[Reading:Equal] Temperatures: %f (%f); Humidity: %f (%f)", r.Temperature, s.Temperature, r.Humidity, s.Humidity)
	if r.Temperature != s.Temperature {
		log.Print("[Reading:Equal] Temperatures don't match")
	}
	if r.Humidity != s.Humidity {
		log.Print("[Reading:Equal] Humidities don't match")
	}
	return r.Temperature == s.Temperature && r.Humidity == s.Humidity
}

// ToString converts a Reading to a string
func (r *Reading) String() string {
	return fmt.Sprintf("Temperature: %.04f; Humidity: %.04f", r.Temperature, r.Humidity)
}

// NewReading returns a new Reading
func NewReading(t, h float64) *Reading {
	return &Reading{
		Temperature: t,
		Humidity:    h,
	}
}

// Unmarshall converts an encoded reading into a Reading
func Unmarshall(req []byte) (*Reading, error) {
	// 00 01 02 03 04
	// T2 T1 HX ?? ??
	l := len(req)
	if l != 5 {
		log.Printf("[X] Expecting 5 bytes; got %d", l)
		return &Reading{}, fmt.Errorf("Expecting 5 bytes got %d", l)
	}
	// Temperature is stored little endian
	t := float64(int(binary.LittleEndian.Uint16(req[0:2]))) / 100.0
	h := float64(req[2]) / 100.0
	return &Reading{
		Temperature: t,
		Humidity:    h,
	}, nil
}
