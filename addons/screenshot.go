package addons

import (
	"os/exec"

	"github.com/magicmonkey/go-streamdeck"
	sdactionhandlers "github.com/magicmonkey/go-streamdeck/actionhandlers"
	"github.com/magicmonkey/go-streamdeck/buttons"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Screenshot struct {
	SD *streamdeck.StreamDeck
}

func (s *Screenshot) Init() {
	// actually don't need to do anything for this action type
}

func (s *Screenshot) Buttons() {
	// Command
	shotbutton, _ := buttons.NewImageFileButton(viper.GetString("buttons.images") + "/screenshot.png")
	shotaction := &sdactionhandlers.CustomAction{}
	shotaction.SetHandler(func(btn streamdeck.Button) {
		go takeScreenshot()
	})
	shotbutton.SetActionHandler(shotaction)
	s.SD.AddButton(15, shotbutton)
}

func takeScreenshot() {
	log.Debug().Msg("Taking screenshot with delay...")
	cmd := exec.Command("/usr/bin/gnome-screenshot", "-w", "-d", "2")
	if err := cmd.Run(); err != nil {
		log.Warn().Err(err)
	}
	log.Debug().Msg("Taken screenshot")
}
