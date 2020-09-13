package addons

import (
	"encoding/json"
	"image/color"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/buttons"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type MqttThing struct {
	SD          *streamdeck.StreamDeck
	mqtt_client mqtt.Client
}

func (p *MqttThing) Init() {
	// Initialise MQTT to use the shelf light features
	log.Debug().Msg("Connecting to MQTT...")
	opts := mqtt.NewClientOptions().AddBroker("tcp://10.1.0.1:1883").SetClientID("go-streamdeck")
	p.mqtt_client = mqtt.NewClient(opts)
	if conn_token := p.mqtt_client.Connect(); conn_token.Wait() && conn_token.Error() != nil {
		log.Warn().Err(conn_token.Error()).Msg("Cannot connect to MQTT")
	}
}

// Set up buttons
type LEDColour struct {
	Red   uint8 `mapstructure:"red" json:"red"`
	Green uint8 `mapstructure:"green" json:"green"`
	Blue  uint8 `mapstructure:"blue" json:"blue"`
	White uint8 `mapstructure:"white" json:"white"`
}

type PlugDevice struct {
	Name     string `mapstructure:"name"`
	ButtonId int    `mapstructure:"button"`
	ImageOn  string `mapstructure:"image_on"`
	ImageOff string `mapstructure:"image_off"`
}

type ManyLights struct {
	Colour1 LEDColour `mapstructure:"colour1"`
	Colour2 LEDColour `mapstructure:"colour2"`
}

func (p *MqttThing) Buttons() {
	var lights []ManyLights
	viper.UnmarshalKey("shelf_lights", &lights)
	button_index := 18

	for _, light := range lights {
		colour := color.RGBA{light.Colour1.Red, light.Colour1.Green, light.Colour1.Blue, 255}
		lbutton := buttons.NewColourButton(colour)
		lbutton.SetActionHandler(&MQTTAction{Colour: light.Colour1, Colour2: light.Colour2, Client: p.mqtt_client})
		p.SD.AddButton(button_index, lbutton)
		button_index = button_index + 1
	}

	// on/off plugs
	var buttons_plug map[string]PlugDevice // MQTT-enabled on/off plugs
	viper.UnmarshalKey("plug_devices", &buttons_plug)
	for device, deets := range buttons_plug {
		// assume off, we can't get state
		image := viper.GetString("buttons.images") + "/" + deets.ImageOff
		plugbutton, err := buttons.NewImageFileButton(image)
		if err == nil {
			plugaction := &PlugAction{Client: p.mqtt_client, Device: device, State: 0, ImageOn: deets.ImageOn, ImageOff: deets.ImageOff}
			plugbutton.SetActionHandler(plugaction)
			p.SD.AddButton(deets.ButtonId, plugbutton)
		} else {
			log.Warn().Err(err)
		}
	}
}

// Button action handler
type PlugAction struct {
	Client   mqtt.Client
	Device   string
	State    int
	ImageOn  string
	ImageOff string
}

func (action *PlugAction) Pressed(btn streamdeck.Button) {
	imagebutton := btn.(*buttons.ImageFileButton)

	// Toggle! Are we on? Turn off! Not on? Turn on!
	if action.State == 1 {
		// on! Turn off
		token := action.Client.Publish("/house/plug/"+action.Device, 0, false, "0")
		token.Wait()
		imagebutton.SetFilePath(viper.GetString("buttons.images") + "/" + action.ImageOff)
		action.State = 0
	} else {
		// off! So turn on
		token := action.Client.Publish("/house/plug/"+action.Device, 0, false, "1")
		token.Wait()
		imagebutton.SetFilePath(viper.GetString("buttons.images") + "/" + action.ImageOn)
		action.State = 1
	}
}

// Lights Action handler
type MQTTAction struct {
	Client  mqtt.Client
	Colour  LEDColour
	Colour2 LEDColour
	btn     streamdeck.Button
}

func (action *MQTTAction) Pressed(btn streamdeck.Button) {
	var targetColour []LEDColour
	targetColour = append(targetColour, action.Colour)
	targetColour = append(targetColour, action.Colour2)

	payload, _ := json.Marshal(targetColour)
	log.Debug().Msg(string(payload))
	token := action.Client.Publish("/shelf/lights", 0, false, payload)
	token.Wait()
}
