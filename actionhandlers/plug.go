package actionhandlers

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	streamdeck "github.com/magicmonkey/go-streamdeck"
	buttons "github.com/magicmonkey/go-streamdeck/buttons"
	"github.com/spf13/viper"
)

type PlugAction struct {
	Client   mqtt.Client
	State    int
	ImageOn  string
	ImageOff string
}

func (action *PlugAction) Pressed(btn streamdeck.Button) {
	imagebutton := btn.(*buttons.ImageFileButton)

	// Toggle! Are we on? Turn off! Not on? Turn on!
	if action.State == 1 {
		// on! Turn off
		token := action.Client.Publish("/house/plug/shelf", 0, false, "0")
		token.Wait()
		imagebutton.SetFilePath(viper.GetString("buttons.images") + "/" + action.ImageOff)
		action.State = 0
	} else {
		// off! So turn on
		token := action.Client.Publish("/house/plug/shelf", 0, false, "1")
		token.Wait()
		imagebutton.SetFilePath(viper.GetString("buttons.images") + "/" + action.ImageOn)
		action.State = 1
	}

}
