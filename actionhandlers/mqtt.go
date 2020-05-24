package actionhandlers

import (
	"encoding/json"
	"image/color"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	streamdeck "github.com/magicmonkey/go-streamdeck"
	"github.com/rs/zerolog/log"
)

type MQTTAction struct {
	Client mqtt.Client
	Colour color.RGBA
	btn    streamdeck.Button
}

type Colour struct {
	Red   uint8 `json:"red"`
	Green uint8 `json:"green"`
	Blue  uint8 `json:"blue"`
}

func (action *MQTTAction) Pressed(btn streamdeck.Button) {

	targetColour := Colour{Red: action.Colour.R, Green: action.Colour.G, Blue: action.Colour.B}

	payload, _ := json.Marshal(targetColour)
	log.Debug().Msg(string(payload))
	token := action.Client.Publish("/shelf/lights", 0, false, payload)
	token.Wait()
}
