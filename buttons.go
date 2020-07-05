package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	obsws "github.com/christopher-dG/go-obs-websocket"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/lornajane/streamdeck-tricks/actionhandlers"
	"github.com/magicmonkey/go-streamdeck"
	sdactionhandlers "github.com/magicmonkey/go-streamdeck/actionhandlers"
	buttons "github.com/magicmonkey/go-streamdeck/buttons"
	sddecorators "github.com/magicmonkey/go-streamdeck/decorators"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	_ "github.com/godbus/dbus"
	"github.com/sqp/pulseaudio"
)

var mqtt_client mqtt.Client
var obs_client obsws.Client
var obs_current_scene string
var pulse *pulseaudio.Client

type PlugDevice struct {
	Name     string `mapstructure:"name"`
	ButtonId int    `mapstructure:"button"`
	ImageOn  string `mapstructure:"image_on"`
	ImageOff string `mapstructure:"image_off"`
}

type ObsScene struct {
	Name     string `mapstructure:"name"`
	Image    string `mapstructure:"image"`
	ButtonId int
}

func (scene *ObsScene) SetButtonId(id int) {
	scene.ButtonId = id
}

var buttons_plug map[string]PlugDevice // MQTT-enabled on/off plugs
var buttons_obs map[string]*ObsScene   // scene name and image name

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
			scene := strings.ToLower(e.(obsws.SwitchScenesEvent).SceneName)
			log.Info().Msg("Old scene: " + obs_current_scene)
			// undecorate the old
			if oldb, ok := buttons_obs[obs_current_scene]; ok {
				log.Info().Int("button", oldb.ButtonId).Msg("Clear original button decoration")
				sd.UnsetDecorator(oldb.ButtonId)
			}
			// decorate the new
			log.Info().Msg("New scene: " + scene)
			if eventb, ok := buttons_obs[scene]; ok {
				decorator2 := sddecorators.NewBorder(5, color.RGBA{255, 0, 0, 255})
				log.Info().Int("button", eventb.ButtonId).Msg("Highlight new scene button")
				sd.SetDecorator(eventb.ButtonId, decorator2)
			}
			obs_current_scene = scene
		})
	}

	// Get some Audio Setup
	pulse = getPulseConnection()

	// on/off plugs
	viper.UnmarshalKey("plug_devices", &buttons_plug)
	for device, deets := range buttons_plug {
		fmt.Println(deets.Name)
		// assume off, we can't get state
		image := viper.GetString("buttons.images") + "/" + deets.ImageOff
		plugbutton, err := buttons.NewImageFileButton(image)
		if err == nil {
			plugaction := &actionhandlers.PlugAction{Client: mqtt_client, Device: device, State: 0, ImageOn: deets.ImageOn, ImageOff: deets.ImageOff}
			plugbutton.SetActionHandler(plugaction)
			sd.AddButton(deets.ButtonId, plugbutton)
		} else {
			log.Warn().Err(err)
		}
	}

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

	// OBS Scenes to Buttons
	buttons_obs = make(map[string]*ObsScene)
	viper.UnmarshalKey("obs_scenes", &buttons_obs)

	if obs_client.Connected() == true {
		// offset for what number button to start at
		offset := 0
		image_path := viper.GetString("buttons.images")
		var image string

		// what scenes do we have? (max 8 for the top row of buttons)
		scene_req := obsws.NewGetSceneListRequest()
		scenes, err := scene_req.SendReceive(obs_client)
		if err != nil {
			log.Warn().Err(err)
		}
		obs_current_scene = strings.ToLower(scenes.CurrentScene)

		// make buttons for these scenes
		for i, scene := range scenes.Scenes {
			log.Debug().Msg("Scene: " + scene.Name)
			image = ""
			oaction := &actionhandlers.OBSSceneAction{Scene: scene.Name, Client: obs_client}
			sceneName := strings.ToLower(scene.Name)

			if s, ok := buttons_obs[sceneName]; ok {
				if s.Image != "" {
					image = image_path + s.Image
				}
			} else {
				// there wasn't an entry in the buttons for this scene so add one
				buttons_obs[sceneName] = &ObsScene{}
			}

			if image != "" {
				// try to make an image button

				obutton, err := buttons.NewImageFileButton(image)
				if err == nil {
					obutton.SetActionHandler(oaction)
					sd.AddButton(i+offset, obutton)
					// store which button we just set
					buttons_obs[sceneName].SetButtonId(i + offset)
				} else {
					// something went wrong with the image, use a default one
					image = image_path + "/play.jpg"
					obutton, err := buttons.NewImageFileButton(image)
					if err == nil {
						obutton.SetActionHandler(oaction)
						sd.AddButton(i+offset, obutton)
						// store which button we just set
						buttons_obs[sceneName].SetButtonId(i + offset)
					}
				}
			} else {
				// use a text button
				oopbutton := buttons.NewTextButton(scene.Name)
				oopbutton.SetActionHandler(oaction)
				sd.AddButton(i+offset, oopbutton)
				// store which button we just set
				buttons_obs[sceneName].SetButtonId(i + offset)
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

	// Sounds (bottom row, button 24+)
	sound_offset := 24
	for i := 0; i < 4; i++ {
		audiobutton := buttons.NewTextButton(strconv.Itoa(i + 1))
		audiobutton.SetActionHandler(&actionhandlers.OSCAction{Track: i + 1})
		sd.AddButton(sound_offset+i, audiobutton)
	}
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
