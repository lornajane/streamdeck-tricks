package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"os/exec"
	"strconv"
	"time"

	obsws "github.com/christopher-dG/go-obs-websocket"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/lornajane/streamdeck-tricks/actionhandlers"
	"github.com/magicmonkey/go-streamdeck"
	sdactionhandlers "github.com/magicmonkey/go-streamdeck/actionhandlers"
	buttons "github.com/magicmonkey/go-streamdeck/buttons"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	_ "github.com/godbus/dbus"
	belkin "github.com/magicmonkey/gobelkinwemo"
	"github.com/sqp/pulseaudio"
)

var mqtt_client mqtt.Client
var obs_client obsws.Client
var obs_current_scene string
var pulse *pulseaudio.Client

type Wemo struct {
	Name     string `mapstructure:"name"`
	ButtonId int    `mapstructure:"button"`
	ImageOn  string `mapstructure:"image_on"`
	ImageOff string `mapstructure:"image_off"`
}

var buttons_wemo map[string]Wemo  // button ID and Wemo device configs
var buttons_obs map[string]string // scene name and image name

type LEDColour struct {
	Red   uint8 `mapstructure:"red"`
	Green uint8 `mapstructure:"green"`
	Blue  uint8 `mapstructure:"blue"`
}

// InitButtons sets up initial button prompts
func InitButtons() {
	// Initialise MQTT to use the shelf light features
	mqtt_client = connectMQTT()

	// Initialise OBS to use OBS features (requires websockets plugin in OBS)
	obs_client = connectOBS()

	if obs_client.Connected() == true {
		obs_client.AddEventHandler("SwitchScenes", func(e obsws.Event) {
			// Make sure to assert the actual event type.
			log.Info().Msg("new scene: " + e.(obsws.SwitchScenesEvent).SceneName)
		})
	}

	// Get some Audio Setup
	pulse = getPulseConnection()

	// WEMO plugs
	viper.UnmarshalKey("wemo_devices", &buttons_wemo)
	go startWemoScan()

	// shelf lights
	var lights []LEDColour
	viper.UnmarshalKey("shelf_lights", &lights)
	button_index := 8

	for _, light := range lights {
		colour := color.RGBA{light.Red, light.Green, light.Blue, 255}
		lbutton := buttons.NewColourButton(colour)
		lbutton.SetActionHandler(&actionhandlers.MQTTAction{Colour: colour, Client: mqtt_client})
		sd.AddButton(button_index, lbutton)
		button_index = button_index + 1
	}

	// OBS (this should come from config)
	buttons_obs = make(map[string]string)
	buttons_obs["Camera"] = "/camera.png"
	buttons_obs["Screenshare"] = "/screen-and-cam.png"
	buttons_obs["layout-main-solo"] = "/camera.png"
	buttons_obs["layout-code-solo"] = "/screen-and-cam.png"
	buttons_obs["Secrets"] = "/secrets.png"
	buttons_obs["Offline"] = "/offline.png"
	buttons_obs["layout-offline"] = "/offline.png"
	buttons_obs["layout-starting"] = "/soon.png"
	buttons_obs["layout-main"] = "/copresenters.png"
	buttons_obs["layout-code-remoter"] = "/their-screen.png"
	buttons_obs["layout-secret"] = "/secrets.png"
	buttons_obs["layout-android"] = "/android-and-cam.png"

	if obs_client.Connected() == true {
		// offset for what number button to start at
		offset := 0
		image_path := viper.GetString("buttons.images")
		var image string

		// what scenes do we have? (max 8)
		scene_req := obsws.NewGetSceneListRequest()
		scenes, err := scene_req.SendReceive(obs_client)
		if err != nil {
			log.Warn().Err(err)
		}
		// fmt.Printf("%#v\n", scenes.CurrentScene)
		// fmt.Printf("%#v\n", scenes.Scenes[2])

		// make buttons for these scenes
		for i, scene := range scenes.Scenes {
			log.Debug().Msg("Scene: " + scene.Name)
			image = ""

			if buttons_obs[scene.Name] != "" {
				image = image_path + buttons_obs[scene.Name]
			}

			oaction := &actionhandlers.OBSSceneAction{Scene: scene.Name, Client: obs_client}
			if image != "" {
				// try to make an image button
				obutton, err := buttons.NewImageFileButton(image)
				if err == nil {
					obutton.SetActionHandler(oaction)
					sd.AddButton(i+offset, obutton)
				} else {
					image = image_path + "/play.jpg"
					obutton, err := buttons.NewImageFileButton(image)
					if err == nil {
						obutton.SetActionHandler(oaction)
						sd.AddButton(i+offset, obutton)
					}
				}
			} else {
				// use a text button
				oopbutton := buttons.NewTextButton(scene.Name)
				oopbutton.SetActionHandler(oaction)
				sd.AddButton(i+offset, oopbutton)
			}

			// only need a few scenes
			if i > 6 {
				break
			}
		}
	}

	// Command
	shotbutton, _ := buttons.NewImageFileButton(viper.GetString("buttons.images") + "/screenshot.png")
	shotaction := &sdactionhandlers.CustomAction{}
	shotaction.SetHandler(func(btn streamdeck.Button) {
		go takeScreenshot()
	})
	shotbutton.SetActionHandler(shotaction)
	sd.AddButton(15, shotbutton)
}

func takeScreenshot() {
	log.Debug().Msg("Taking screenshot with delay...")
	cmd := exec.Command("/usr/bin/gnome-screenshot", "-w", "-d", "2")
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Run(); err != nil {
		log.Warn().Err(err)
	}

	slurp, _ := ioutil.ReadAll(stderr)
	fmt.Printf("%s\n", slurp)
	slurp2, _ := ioutil.ReadAll(stdout)
	fmt.Printf("%s\n", slurp2)

	log.Debug().Msg("Taken screenshot")
}

func connectMQTT() mqtt.Client {
	log.Debug().Msg("Connecting to MQTT...")
	opts := mqtt.NewClientOptions().AddBroker("tcp://10.1.0.1:1883").SetClientID("go-streamdeck")
	mqtt_client = mqtt.NewClient(opts)
	if conn_token := mqtt_client.Connect(); conn_token.Wait() && conn_token.Error() != nil {
		log.Warn().Err(conn_token.Error()).Msg("Cannot connect to MQTT")
	}
	return mqtt_client
}

func connectOBS() obsws.Client {
	log.Debug().Msg("Connecting to OBS...")
	log.Info().Msgf("%#v\n", viper.Get("obs.host"))
	obs_client = obsws.Client{Host: "localhost", Port: 4444}
	err := obs_client.Connect()
	if err != nil {
		log.Warn().Err(err).Msg("Cannot connect to OBS")
	}
	return obs_client
}

/*
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
	}
}
*/

type AppPulse struct {
	Client *pulseaudio.Client
}

func getPulseConnection() *pulseaudio.Client {
	isLoaded, e := pulseaudio.ModuleIsLoaded()
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
		log.Warn().Err(e).Msg(msg)
	}
}

// Wemo functions from magicmonkey modified library
func startWemoScan() {
	device1, err1 := belkin.NewDeviceFromURL("http://10.1.0.170:49153/setup.xml", 2*time.Second)
	if err1 != nil {
		log.Warn().Msg("Device 170 not found")
	} else {
		gotWemoDevice(*device1)
	}

	device2, err2 := belkin.NewDeviceFromURL("http://10.1.0.117:49153/setup.xml", 2*time.Second)
	if err2 != nil {
		log.Warn().Msg("Device 117 not found")
	} else {
		gotWemoDevice(*device2)
	}
	// the scan is the official approach but wasn't very reliable
	// err := belkin.ScanWithCallback(belkin.DTInsight, 10, gotWemoDevice)
}

func gotWemoDevice(device belkin.Device) {
	device.Load(1 * time.Second)
	state, err := device.FetchBinaryState(1 * time.Second)
	if err != nil {
		log.Warn().Err(err)
	}
	log.Info().Msg("Found device " + device.FriendlyName)
	log.Debug().Msg("Current device state: " + strconv.Itoa(state)) // 0, 1 or 8 (for standby)

	for _, deets := range buttons_wemo {
		log.Debug().Int("button id", deets.ButtonId).Msg(deets.Name)
		if deets.Name == device.FriendlyName {
			image := viper.GetString("buttons.images") + "/" + deets.ImageOff
			if state == 1 {
				image = viper.GetString("buttons.images") + "/" + deets.ImageOn
			}
			wemobutton, err := buttons.NewImageFileButton(image)
			if err == nil {
				wemoaction := &actionhandlers.WemoAction{Device: device, State: device.BinaryState, ImageOn: deets.ImageOn, ImageOff: deets.ImageOff}
				wemobutton.SetActionHandler(wemoaction)
				sd.AddButton(deets.ButtonId, wemobutton)
			} else {
				log.Warn().Err(err)
			}
		}
	}

}
