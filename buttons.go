package main

import (
	"encoding/json"
	"fmt"
	"os/exec"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/magicmonkey/go-streamdeck"
	"github.com/spf13/viper"
)

var button_images [32]string
var mqtt_client mqtt.Client

// InitButtons sets up initial button prompts
func InitButtons(sd *streamdeck.Device) {
	// shelf lights

	sd.WriteColorToButton(255, 0, 255, 8)
	sd.WriteColorToButton(0, 0, 255, 9)
	sd.WriteColorToButton(255, 255, 0, 10)

	opts := mqtt.NewClientOptions().AddBroker("tcp://10.1.0.1:1883").SetClientID("go-streamdeck")
	mqtt_client = mqtt.NewClient(opts)
	if conn_token := mqtt_client.Connect(); conn_token.Wait() && conn_token.Error() != nil {
		panic(conn_token.Error())
	}

}

// MyButtonPress reacts to a button being pressed
func MyButtonPress(btnIndex int, sd *streamdeck.Device) {
	switch btnIndex {
	case 7:
		cmd := exec.Command("xeyes")
		cmd.Start()
	case 8:
		SetShelfLights(Color{Red: 200, Blue: 200})
	case 9:
		SetShelfLights(Color{Blue: 200})
	case 10:
		SetShelfLights(Color{Red: 160, Green: 200})
	default:
		ToggleImageOnButton(sd, btnIndex, viper.GetString("images_buttons")+"/play.jpg")
	}
}

// Set an image on a button, keep track of that. If the image is set, send black
func ToggleImageOnButton(sd *streamdeck.Device, btnIndex int, image string) {
	if button_images[btnIndex] == "" {
		sd.WriteImageToButton(image, btnIndex)
		button_images[btnIndex] = image
	} else {
		sd.WriteColorToButton(0, 0, 0, btnIndex)
		button_images[btnIndex] = ""
	}
}

type Color struct {
	Red   int `json:"red"`
	Green int `json:"green"`
	Blue  int `json:"blue"`
}

// SetShelfLights sends MQTT messages to my neopixel-enabled shelf
func SetShelfLights(targetColor Color) {
	payload, _ := json.Marshal(targetColor)
	fmt.Printf("%s\n", string(payload))
	token := mqtt_client.Publish("/shelf/lights", 0, false, payload)
	token.Wait()

}
