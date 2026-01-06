package addons

import (
	"encoding/json"
	"image/color"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
)

type MqttThing struct {
	mqtt_client mqtt.Client
	fade_index  int
	snake_state bool
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

func (p *MqttThing) KeyPress(char rune) {
	log.Debug().Msgf("I got key: %c", char)
  switch char {
  case '0':
      p.TriggerFirework(0)
  case '1':
      p.SetBackground(1, 3, 3)  // cyan-ish
  case '2':
      p.SetBackground(6, 5, 1)  // cyan-ish
  case '3':
      p.SetBackground(6, 1, 4)  // cyan-ish
  case '7':
      p.TriggerFade()
  }
}

// SetBackground sets a solid background color on the LED wall
func (p *MqttThing) SetBackground(red, green, blue uint8) {
	p.ClearPixels()
	options := LEDWallBg{Action: "background", Red: red, Green: green, Blue: blue}
	payload, _ := json.Marshal(options)
	log.Debug().Msg(string(payload))
	token := p.mqtt_client.Publish("/ledwall/1/request", 0, false, payload)
	token.Wait()
}

// TriggerFade cycles through predefined fade patterns
func (p *MqttThing) TriggerFade() {
	fade_options := []LEDWallFade{
		{
			BL: color.RGBA{R: 50, G: 00, B: 00, A: 255},
			TL: color.RGBA{R: 50, G: 50, B: 00, A: 255},
			BR: color.RGBA{R: 00, G: 00, B: 50, A: 255},
			TR: color.RGBA{R: 50, G: 50, B: 00, A: 255},
		},
		{
			BL: color.RGBA{R: 15, G: 00, B: 10, A: 255},
			TL: color.RGBA{R: 20, G: 20, B: 00, A: 255},
			BR: color.RGBA{R: 10, G: 00, B: 10, A: 255},
			TR: color.RGBA{R: 30, G: 00, B: 20, A: 255},
		},
		{
			BL: color.RGBA{R: 0, G: 0, B: 50, A: 255},
			TL: color.RGBA{R: 0, G: 20, B: 30, A: 255},
			BR: color.RGBA{R: 10, G: 0, B: 20, A: 255},
			TR: color.RGBA{R: 50, G: 40, B: 0, A: 255},
		},
		{
			BL: color.RGBA{R: 00, G: 20, B: 10, A: 255},
			TL: color.RGBA{R: 10, G: 10, B: 30, A: 255},
			BR: color.RGBA{R: 20, G: 00, B: 10, A: 255},
			TR: color.RGBA{R: 00, G: 10, B: 15, A: 255},
		},
	}

	pixels := Fade(fade_options[p.fade_index])
	for pixnum, pix := range pixels {
		payload, _ := json.Marshal(LEDWallPixel{Action: "pixel", Num: uint8(pixnum), Red: pix.R, Green: pix.G, Blue: pix.B})
		token := p.mqtt_client.Publish("/ledwall/1/request", 0, false, payload)
		token.Wait()
	}

	p.fade_index = p.fade_index + 1
	if p.fade_index >= len(fade_options) {
		p.fade_index = 0
	}
}

// TriggerFirework triggers a firework effect with optional hue
func (p *MqttThing) TriggerFirework(hue uint8) {
	options := LEDWallFirework{Action: "firework", Hue: hue}
	payload, _ := json.Marshal(options)
	log.Debug().Msg(string(payload))
	token := p.mqtt_client.Publish("/ledwall/1/request", 0, false, payload)
	token.Wait()
}

// ToggleSnake toggles the snake effect on/off
func (p *MqttThing) ToggleSnake() {
	options := LEDWallSnake{Action: "snake", Enabled: p.snake_state}
	payload, _ := json.Marshal(options)
	log.Debug().Msg(string(payload))
	token := p.mqtt_client.Publish("/ledwall/1/request", 0, false, payload)
	token.Wait()

	// toggle state for next press
	p.snake_state = !p.snake_state
}

// ClearPixels clears all pixels on the LED wall
func (p *MqttThing) ClearPixels() {
	payload := "{\"action\": \"clearpixels\"}"
	log.Debug().Msg(string(payload))
	token := p.mqtt_client.Publish("/ledwall/1/request", 0, false, payload)
	token.Wait()
}

type MQTTFadeAction struct {
	Client      mqtt.Client
	fade_index  int
}

type MQTTBgAction struct {
	Client  mqtt.Client
	Options LEDWallBg
}

type MQTTClearAction struct {
	Client  mqtt.Client
}


type MQTTFireworkAction struct {
	Client  mqtt.Client
	Options LEDWallFirework
}

type MQTTSnakeAction struct {
	Client      mqtt.Client
	Options     LEDWallSnake
	snake_state bool
}

