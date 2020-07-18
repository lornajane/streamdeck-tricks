package addons

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/buttons"
	sdactionhandlers "github.com/magicmonkey/go-streamdeck/actionhandlers"
	"github.com/spf13/viper"
	"github.com/rs/zerolog/log"
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
