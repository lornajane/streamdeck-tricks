package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os/exec"

	"github.com/christopher-dG/go-obs-websocket"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/magicmonkey/go-streamdeck"
	"github.com/spf13/viper"
)

var button_images [32]string
var mqtt_client mqtt.Client
var obs_client obsws.Client

// InitButtons sets up initial button prompts
func InitButtons(sd *streamdeck.Device) {
	// shelf lights
	sd.WriteColorToButton(8, color.RGBA{255, 0, 255, 255})
	sd.WriteColorToButton(9, color.RGBA{0, 0, 255, 255})
	sd.WriteColorToButton(10, color.RGBA{255, 255, 0, 255})

	ToggleImageOnButton(sd, 24, viper.GetString("images_buttons")+"/flamingo.jpg")
	ToggleImageOnButton(sd, 25, viper.GetString("images_buttons")+"/flamingo.jpg")

	// Initialise MQTT to use the shelf light features
	opts := mqtt.NewClientOptions().AddBroker("tcp://10.1.0.1:1883").SetClientID("go-streamdeck")
	mqtt_client = mqtt.NewClient(opts)
	if conn_token := mqtt_client.Connect(); conn_token.Wait() && conn_token.Error() != nil {
		panic(conn_token.Error())
	}

	// Initialise OBS to use OBS features (requires websockets plugin in OBS)
	obs_client = obsws.Client{Host: "localhost", Port: 4444}
	if err := obs_client.Connect(); err != nil {
		log.Fatal(err)
	}

	// probably need to react to scene changes, log it for now
	obs_client.AddEventHandler("SwitchScenes", func(e obsws.Event) {
		// Make sure to assert the actual event type.
		log.Println("new scene:", e.(obsws.SwitchScenesEvent).SceneName)
	})

}

// MyButtonPress reacts to a button being pressed
func MyButtonPress(btnIndex int, sd *streamdeck.Device) {
	switch btnIndex {
	case 15:
		cmd := exec.Command("xeyes")
		cmd.Start()
	case 8:
		SetShelfLights(Color{Red: 200, Blue: 200})
	case 9:
		SetShelfLights(Color{Blue: 200})
	case 10:
		SetShelfLights(Color{Red: 160, Green: 200})
	case 24:
		fmt.Println("Set scene: Secrets")
		req := obsws.NewSetCurrentSceneRequest("Secrets")
		resp, err := req.SendReceive(obs_client)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%#v\n", resp)

	case 25:
		fmt.Println("Set scene: Soon")
		req := obsws.NewSetCurrentSceneRequest("Soon")
		resp, err := req.SendReceive(obs_client)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%#v\n", resp)

	default:
		ToggleImageOnButton(sd, btnIndex, viper.GetString("images_buttons")+"/play.jpg")
	}
}

// Set an image on a button, keep track of that. If the image is set, send black
func ToggleImageOnButton(sd *streamdeck.Device, btnIndex int, image string) {
	fmt.Println(image)
	if button_images[btnIndex] == "" {
		sd.WriteImageToButton(btnIndex, image)
		button_images[btnIndex] = image
	} else {
		sd.WriteColorToButton(btnIndex, color.RGBA{0, 0, 0, 255})
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
