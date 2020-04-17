# Golang Xiaomi Bluetooth Temperature|Humidity (LYWSD03MMC) 2nd Gen

## Credits:

+ [Johnathan McDowell](https://github.com/u1f35c) for [gomijia](https://github.com/u1f35c/gomijia) and [help](https://github.com/u1f35c/gomijia/issues/1)
+ [JsBergbau](https://github.com/JsBergbau) for [help](https://github.com/JsBergbau/MiTemperature2/issues/29#issuecomment-614314939)

## Summary

The [LYWSD03MMC](https://www.google.com/search?q=LYWSD03MMC) is a second-generation Xiaomi Thermometer|Humidity sensor. The device is BLE-enabled and (supposedly) works in conjunction with the Mi Home app (though I've been unsuccessfully in connecting my devices with the app). I plan to connect the (4) devices that I bought (for $20) to ESPs running the most-excellent [ESPHome](https://esphome.io) but, until new ESP devices arrive, I was interested in connecting to the devices using Golang.

I found Jonathan's [gomijia](https://github.com/u1f35c/gomijia) but this does not work with these second-generation devices. Thanks to Jonatha for his code and his help. After much hacking around, I learned two key facts:

+ The LYWSD03MMC devices don't advertise the temperature & humidity; instead a connection must be made to the device
+ Somehow (!) clever folks (see JsBergbau's working Python [code](https://github.com/JsBergbau/MiTemperature2)) worked out how to receive notifications from the devices (see below)

## Running

You'll need:

+ At least 1 LYWSD03MMC and you'll want the MAC addresses for your devices.
+ An MQTT broker to publish temperature|humidity readings

Create a configuration file and enumerate your devices and their MAC addresses:

```ini

[MQTT]
host=localhost
port=1883
#user=user
#pass=pass

[Devices]
device1=a4:c1:38:11:11:11
device2=a4:c1:38:22:22:22
...
```

You'll need to run the code as root because of the BLE connecting:

```bash
${which go} run github.com/DazWilkin/gomijia2 --config_file=/path/to/config.ini
```

You'll receive lots of logging output but you should see:

```
[main] Reading configuration
[Config] Loading Configuration (./config.ini)
[Config] Device 00: xiaomi1 (a4:c1:38:00:00:00)
[Config] Device 01: xiaomi2 (a4:c1:38:00:00:00)
[Config] Device 02: xiaomi3 (a4:c1:38:00:00:00)
[Config] Device 03: xiaomi4 (a4:c1:38:00:00:00)
[main] Starting Linux Device
[main] MQTT broker: tcp://lcoalhost:1883
[MQTT:Connect] Connecting: xiaomi
[MQTT] Creating
[MQTT] Connecting
[main:xiaomi1] Dialing (a4:c1:38:3c:7c:e7)
[main:xiaomi1] Registering handler
[main:xiaomi1] (0035) Registering Temperature|Humidity Handler
[main:xiaomi2] Dialing (a4:c1:38:00:00:00)
...
[handler:xiaomi1] Temperature: 22.0300; Humidity: 0.4200 (9b082ab90b)
[handler:xiaomi2] Temperature: 21.5600; Humidity: 0.4200 (6c082a600b)
[handler:xiaomi4] Temperature: 22.2900; Humidity: 0.4100 (b50829dd0b)
[handler:xiaomi3] Temperature: 22.0900; Humidity: 0.4200 (a1082ada0a)
[handler:xiaomi1] Temperature: 22.0000; Humidity: 0.4200 (98082ab90b)
[handler:xiaomi2] Temperature: 21.6000; Humidity: 0.4200 (70082a600b)
[handler:xiaomi4] Temperature: 22.3100; Humidity: 0.4100 (b70829dd0b)
```

> **NB** I've removed some entries from the above for clarity

I've configured Inuits [MQTTGateway](https://github.com/inuits/mqttgateway) to subscribe to topics and expose these as Prometheus metrics.

I've configured Prometheus to scrape the MQTTGateway metrics.

Here's the result:


