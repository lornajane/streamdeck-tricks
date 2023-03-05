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

type LEDWallBg struct {
	Action string `json:"action"`
	Red    uint8  `json:"r"`
	Green  uint8  `json:"g"`
	Blue   uint8  `json:"b"`
}

type LEDWallPixel struct {
	Action string `json:"action"`
	Num    uint8  `json:"num"`
	Red    uint8  `json:"r"`
	Green  uint8  `json:"g"`
	Blue   uint8  `json:"b"`
}

type LEDWallFirework struct {
	Action string `json:"action"`
	Hue    uint8  `json:"b,omitempty"`
}
type LEDWallSnake struct {
	Action  string `json:"action"`
	Enabled bool   `json:"enabled"`
}

type LEDWallFade struct {
	BL, TL, BR, TR color.RGBA
}

// Set up buttons
func (p *MqttThing) Buttons() {
	button_index := 18

	bgoptions3 := LEDWallBg{"background", 1, 3, 3}
	lbutton3 := buttons.NewColourButton(color.RGBA{155, 255, 255, 255})
	lbutton3.SetActionHandler(&MQTTBgAction{Options: bgoptions3, Client: p.mqtt_client})
	p.SD.AddButton(button_index, lbutton3)
	button_index = button_index + 1

	bgoptions4 := LEDWallBg{"background", 6, 1, 4}
	lbutton4 := buttons.NewColourButton(color.RGBA{255, 155, 155, 255})
	lbutton4.SetActionHandler(&MQTTBgAction{Options: bgoptions4, Client: p.mqtt_client})
	p.SD.AddButton(button_index, lbutton4)
	button_index = button_index + 1

	bgoptions2 := LEDWallFade{
		BL: color.RGBA{R: 0, G: 0, B: 0, A: 0},
		TL: color.RGBA{R: 50, G: 0, B: 0, A: 0},
		BR: color.RGBA{R: 0, G: 50, B: 0, A: 0},
		TR: color.RGBA{R: 50, G: 50, B: 0, A: 0},
	}
	lbutton2 := buttons.NewColourButton(color.RGBA{0, 0, 155, 255})
	lbutton2.SetActionHandler(&MQTTFadeAction{Options: bgoptions2, Client: p.mqtt_client})
	p.SD.AddButton(button_index, lbutton2)
	button_index = button_index + 1

	/*
		bgoptions2 := LEDWallBg{"background", 6, 5, 1}
		lbutton2 := buttons.NewColourButton(color.RGBA{255, 255, 155, 255})
		lbutton2.SetActionHandler(&MQTTBgAction{Options: bgoptions2, Client: p.mqtt_client})
		p.SD.AddButton(button_index, lbutton2)
		button_index = button_index + 1
	*/

	bgoptions5 := LEDWallBg{"background", 4, 2, 4}
	lbutton5 := buttons.NewColourButton(color.RGBA{255, 200, 255, 255})
	lbutton5.SetActionHandler(&MQTTBgAction{Options: bgoptions5, Client: p.mqtt_client})
	p.SD.AddButton(button_index, lbutton5)

	fireworkoptions := LEDWallFirework{"firework", 0}
	fbutton, ferr := buttons.NewImageFileButton(viper.GetString("buttons.images") + "/firework-sparkler.png")
	if ferr != nil {
		panic(ferr)
	}
	fbutton.SetActionHandler(&MQTTFireworkAction{Options: fireworkoptions, Client: p.mqtt_client})
	p.SD.AddButton(16, fbutton)

	snakeoptions := LEDWallSnake{"snake", false}
	sbutton := buttons.NewTextButton("Snake")
	sbutton.SetActionHandler(&MQTTSnakeAction{Options: snakeoptions, Client: p.mqtt_client})
	p.SD.AddButton(17, sbutton)

}

type MQTTFadeAction struct {
	Client  mqtt.Client
	Options LEDWallFade
	btn     streamdeck.Button
}

func (action *MQTTFadeAction) Pressed(btn streamdeck.Button) {
	pixels := Fade(action.Options)
	for pixnum, pix := range pixels {
		payload, _ := json.Marshal(LEDWallPixel{Action: "pixel", Num: uint8(pixnum), Red: pix.R, Green: pix.G, Blue: pix.B})
		token := action.Client.Publish("/ledwall/1/request", 0, false, payload)
		token.Wait()
	}
	/*
		payload, _ := json.Marshal(action.Options)
		token := action.Client.Publish("/ledwall/1/request", 0, false, payload)
		token.Wait()
	*/
}

type MQTTBgAction struct {
	Client  mqtt.Client
	Options LEDWallBg
	btn     streamdeck.Button
}

func (action *MQTTBgAction) Pressed(btn streamdeck.Button) {
	payload, _ := json.Marshal(action.Options)
	log.Debug().Msg(string(payload))
	token := action.Client.Publish("/ledwall/1/request", 0, false, payload)
	token.Wait()
}

type MQTTFireworkAction struct {
	Client  mqtt.Client
	Options LEDWallFirework
	btn     streamdeck.Button
}

func (action *MQTTFireworkAction) Pressed(btn streamdeck.Button) {
	payload, _ := json.Marshal(action.Options)
	log.Debug().Msg(string(payload))
	token := action.Client.Publish("/ledwall/1/request", 0, false, payload)
	token.Wait()
}

type MQTTSnakeAction struct {
	Client      mqtt.Client
	Options     LEDWallSnake
	btn         streamdeck.Button
	snake_state bool
}

func (action *MQTTSnakeAction) Pressed(btn streamdeck.Button) {
	// use current snake state before sending, then switch state after
	action.Options.Enabled = action.snake_state
	payload, _ := json.Marshal(action.Options)
	log.Debug().Msg(string(payload))
	token := action.Client.Publish("/ledwall/1/request", 0, false, payload)
	token.Wait()

	// now toggle state ready for next press
	if action.snake_state {
		action.snake_state = false
	} else {
		action.snake_state = true
	}
}
