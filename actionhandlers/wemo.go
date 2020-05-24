package actionhandlers

import (
	"time"

	streamdeck "github.com/magicmonkey/go-streamdeck"
	buttons "github.com/magicmonkey/go-streamdeck/buttons"
	belkin "github.com/magicmonkey/gobelkinwemo"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type WemoAction struct {
	Device   belkin.Device
	State    int
	ImageOn  string
	ImageOff string
}

func (action *WemoAction) Pressed(btn streamdeck.Button) {
	imagebutton := btn.(*buttons.ImageFileButton)
	log.Debug().Msg(action.Device.FriendlyName)

	// Toggle! Are we on? Turn off! Not on? Turn on!
	if action.State == 1 {
		// on! Turn off
		action.Device.TurnOff(1 * time.Second)
		imagebutton.SetFilePath(viper.GetString("buttons.images") + "/" + action.ImageOff)
		action.State = 0
	} else {
		// off! So turn on
		action.Device.TurnOn(1 * time.Second)
		imagebutton.SetFilePath(viper.GetString("buttons.images") + "/" + action.ImageOn)
		action.State = 1
	}

}
