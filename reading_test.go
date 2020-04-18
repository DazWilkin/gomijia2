package main

import (
	"fmt"
	"testing"
)

var (
	tests = map[Reading][]byte{
		{
			Temperature: 22.20,
			Humidity:    0.31,
		}: {0xac, 0x08, 0x1f, 0x33, 0x0b},
		{
			Temperature: 21.92,
			Humidity:    0.31,
		}: {0x90, 0x08, 0x1f, 0xc1, 0x0b},
		{
			Temperature: 15.43,
			Humidity:    0.56,
		}: {0x07, 0x06, 0x38, 0xb6, 0x0a},
	}
)

// TestUnmarshall tests Unmarshall
func TestUnmarshall(t *testing.T) {
	for r, b := range tests {
		t.Run(fmt.Sprintf("%x", b), func(t *testing.T) {
			want := r
			got, err := Unmarshall(b)
			if err != nil {
				t.Errorf("Error")
			}
			if !got.Equal(want) {
				t.Errorf("got: %s, want: %s", got.String(), want.String())
			}
		})
	}
}
