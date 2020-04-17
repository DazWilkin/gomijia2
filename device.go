package main

import (
	"log"
	"time"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/linux"
	"golang.org/x/net/context"
)

var (
	service = ble.MustParse("ebe0ccb0-7a0a-4b0c-8a1a-6ff2997da3a6")
)

var (
	characteristix = map[uint8]ble.UUID{
		35: ble.MustParse("ebe0ccc1-7a0a-4b0c-8a1a-6ff2997da3a6"),
		38: ble.MustParse("00002902-0000-1000-8000-00805f9b34fb"),
	}
)

// Device represents a BLE Device
type Device struct {
	Name   string
	Addr   string
	Client ble.Client
}

// Connect to a Device
func (d *Device) Connect(host *linux.Device) (err error) {
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), 1*time.Minute))
	d.Client, err = host.Dial(ctx, ble.NewAddr(d.Addr))
	return err
}

// Disconnect from a Device
func (d *Device) Disconnect() error {
	return d.Client.CancelConnection()
}

// RegisterHandler registers a Temperature|Humidity handler
func (d *Device) RegisterHandler(mqtt MQTT) {

	// Thanks to other developers
	// Write to handle `0x0038` with value `0x0100` is required to trigger notification of humidity|temperature
	log.Printf("[Device:RegisterHandler:%s] Write to handler", d.Name)
	d.writeToHandler(characteristix[38], []byte{0x01, 0x00})

	log.Printf("[Device:RegisterHandler:%s] Discover services", d.Name)
	services, err := d.Client.DiscoverServices([]ble.UUID{service})
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range services {
		characteristics, err := d.Client.DiscoverCharacteristics([]ble.UUID{
			characteristix[35],
		}, s)
		if err != nil {
			log.Fatal(err)
		}

		for _, c := range characteristics {
			_, err := d.Client.DiscoverDescriptors(nil, c)
			if err != nil {
				log.Print(err)
			}

			// If this Characteristic suports notifications, we'll subscribe to it
			if (c.Property & ble.CharNotify) != 0 {
				if c.CCCD == nil {
					continue
				}
				// Register Handler but only for Characteristic corresponding to handler `0x0035`
				if c.UUID.Equal(characteristix[35]) {
					log.Printf("[main:%s] (%04x) Registering Temperature|Humidity Handler", d.Name, c.Handle)
					// if err := bleClient.Subscribe(c, false, tempHandler(name)); err != nil {
					if err := d.Client.Subscribe(c, false, handlerPublisher(mqtt, d.Name)); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}
}
func (d *Device) writeToHandler(c ble.UUID, b []byte) {
	log.Printf("[Device:writeToHandler:%s] Handler: %s (%x)", d.Name, c.String(), b)
	if p, err := d.Client.DiscoverProfile(true); err != nil {
		if u := p.Find(ble.NewCharacteristic(c)); u != nil {
			if err := d.Client.WriteCharacteristic(u.(*ble.Characteristic), b, false); err != nil {
				log.Print(err)
			}
		}
	}

}
