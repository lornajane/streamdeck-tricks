package actionhandlers

import (
	"image/color"
	"time"

	streamdeck "github.com/magicmonkey/go-streamdeck"
	buttons "github.com/magicmonkey/go-streamdeck/buttons"
	belkin "github.com/magicmonkey/gobelkinwemo"
	"github.com/rs/zerolog/log"
)

type WemoAction struct {
	Device belkin.Device
	State  int
}

func (action *WemoAction) Pressed(btn streamdeck.Button) {
	textbutton := btn.(*buttons.TextButton)
	log.Debug().Msg(action.Device.FriendlyName)

	// Toggle! Are we on? Turn off! Not on? Turn on!
	if action.State == 1 {
		// on! Turn off
		action.Device.TurnOff(1 * time.Second)
		textbutton.SetTextColour(color.RGBA{255, 0, 0, 255})
		action.State = 0
	} else {
		// off! So turn on
		action.Device.TurnOn(1 * time.Second)
		textbutton.SetTextColour(color.RGBA{0, 255, 0, 255})
		action.State = 1
	}

}
