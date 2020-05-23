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

	_ "github.com/godbus/dbus"
	"github.com/sqp/pulseaudio"
)

var button_images [32]string
var mqtt_client mqtt.Client
var obs_client obsws.Client
var pulse *pulseaudio.Client

// InitButtons sets up initial button prompts
func InitButtons(sd *streamdeck.Device) {
	// shelf lights
	sd.WriteColorToButton(8, color.RGBA{255, 0, 255, 255})
	sd.WriteColorToButton(9, color.RGBA{0, 0, 255, 255})
	sd.WriteColorToButton(10, color.RGBA{255, 255, 0, 255})

	ToggleImageOnButton(sd, 24, viper.GetString("images_buttons")+"/flamingo.jpg")
	ToggleImageOnButton(sd, 25, viper.GetString("images_buttons")+"/flamingo.jpg")
	sd.WriteTextToButton(26, "Hi Lorna!", color.RGBA{0, 0, 0, 255}, color.RGBA{0, 255, 255, 255})

	// Initialise MQTT to use the shelf light features
	opts := mqtt.NewClientOptions().AddBroker("tcp://10.1.0.1:1883").SetClientID("go-streamdeck")
	mqtt_client = mqtt.NewClient(opts)
	if conn_token := mqtt_client.Connect(); conn_token.Wait() && conn_token.Error() != nil {
		log.Println(conn_token.Error())
	}

	// Initialise OBS to use OBS features (requires websockets plugin in OBS)
	obs_client = obsws.Client{Host: "localhost", Port: 4444}
	err := obs_client.Connect()
	if err != nil {
		log.Println(err)
	} else {
		// go ahead and set up event handlers

		// probably need to react to scene changes, log it for now
		obs_client.AddEventHandler("SwitchScenes", func(e obsws.Event) {
			// Make sure to assert the actual event type.
			log.Println("new scene:", e.(obsws.SwitchScenesEvent).SceneName)
		})
	}

	// Get some Audio Setup
	pulse = getPulseConnection()

}

// MyButtonPress reacts to a button being pressed
func MyButtonPress(btnIndex int, sd *streamdeck.Device, err error) {
	switch btnIndex {
	case 0:
		sources, _ := pulse.Core().ListPath("Sources")

		for _, src := range sources {
			dev := pulse.Device(src) // Only use the first sink for the test.
			var name string
			var muted bool
			dev.Get("Name", &name)
			dev.Get("Mute", &muted)
			fmt.Println(src, muted, name)

			dev.Set("Mute", true)
		}
	case 15:
		cmd := exec.Command("xeyes")
		cmd.Start()
	case 8:
		SetShelfLights(Color{Red: 200, Blue: 200})
	case 9:
		SetShelfLights(Color{Blue: 200, Green: 30, Red: 50})
	case 10:
		SetShelfLights(Color{Red: 160, Green: 200})
	case 24:
		fmt.Println("Set scene: Camera")
		req := obsws.NewSetCurrentSceneRequest("Camera")
		resp, err := req.SendReceive(obs_client)
		if err != nil {
			fmt.Printf("%#v\n", err)
		}
		fmt.Printf("%#v\n", resp)

	case 25:
		fmt.Println("Set scene: ScreenshareWithCam")
		req := obsws.NewSetCurrentSceneRequest("ScreenshareWithCam")
		resp, err := req.SendReceive(obs_client)
		if err != nil {
			fmt.Printf("%#v\n", err)
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

// SetShelfLights sends MQTT messages to my neopixel-enabled shelf (max value: 200)
func SetShelfLights(targetColor Color) {
	payload, _ := json.Marshal(targetColor)
	fmt.Printf("%s\n", string(payload))
	token := mqtt_client.Publish("/shelf/lights", 0, false, payload)
	token.Wait()
}

type AppPulse struct {
	Client *pulseaudio.Client
}

func getPulseConnection() *pulseaudio.Client {
	isLoaded, e := pulseaudio.ModuleIsLoaded()
	fmt.Printf("%#v\n", e)
	testFatal(e, "test pulse dbus module is loaded")
	if !isLoaded {
		e = pulseaudio.LoadModule()
		testFatal(e, "load pulse dbus module")
	}

	// Connect to the pulseaudio dbus service.
	pulse, e := pulseaudio.New()
	testFatal(e, "connect to the pulse service")
	return pulse
}

func closePulseConnection(pulse *pulseaudio.Client) {
	//defer pulseaudio.UnloadModule()
	defer pulse.Close()
}

func testFatal(e error, msg string) {
	if e != nil {
		log.Fatalln(msg+":", e)
	}
}
