package addons

import (
	"os/exec"

	"github.com/magicmonkey/go-streamdeck"
	sdactionhandlers "github.com/magicmonkey/go-streamdeck/actionhandlers"
	buttons "github.com/magicmonkey/go-streamdeck/buttons"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type WindowMgmt struct {
	SD *streamdeck.StreamDeck
}

func (s *WindowMgmt) Init() {
	// placeholder for initialising windows
}

func (s *WindowMgmt) Buttons() {
	//  Focus OBS
	obsbutton, _ := buttons.NewImageFileButton(viper.GetString("buttons.images") + "/obs-logo.png")
	obsaction := &sdactionhandlers.CustomAction{}
	obsaction.SetHandler(func(btn streamdeck.Button) {
		cmd := exec.Command("/usr/bin/wmctrl", "-a", "OBS ")
		if err := cmd.Run(); err != nil {
			log.Warn().Err(err)
		}
	})
	obsbutton.SetActionHandler(obsaction)
	s.SD.AddButton(14, obsbutton)

	// Focus window called "Stream Manager" for Twitch
	twitchbutton, _ := buttons.NewImageFileButton(viper.GetString("buttons.images") + "/twitch-logo.png")
	twitchaction := &sdactionhandlers.CustomAction{}
	twitchaction.SetHandler(func(btn streamdeck.Button) {
		cmd := exec.Command("/usr/bin/wmctrl", "-a", "Stream Manager")
		if err := cmd.Run(); err != nil {
			log.Warn().Err(err)
		}
	})
	twitchbutton.SetActionHandler(twitchaction)
	s.SD.AddButton(13, twitchbutton)

	// Focus window called "featured.chat"
	chatbutton := buttons.NewTextButton("Chat")
	chataction := &sdactionhandlers.CustomAction{}
	chataction.SetHandler(func(btn streamdeck.Button) {
		cmd := exec.Command("/usr/bin/wmctrl", "-a", "featured.chat")
		if err := cmd.Run(); err != nil {
			log.Warn().Err(err)
		}
	})
	chatbutton.SetActionHandler(chataction)
	s.SD.AddButton(12, chatbutton)

}
