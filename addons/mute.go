package addons

import (
	"bytes"
	"image/color"
	"os/exec"
	"strings"

	"github.com/magicmonkey/go-streamdeck"
	sdactionhandlers "github.com/magicmonkey/go-streamdeck/actionhandlers"
	"github.com/magicmonkey/go-streamdeck/buttons"
	sddecorators "github.com/magicmonkey/go-streamdeck/decorators"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Mute struct {
	SD        *streamdeck.StreamDeck
	Status    bool // true if muted
	Button_id int
}

var source = "source-10992"

func (s *Mute) Init() {
	// not much to initialise but should probably read some config for source name
	// or calculate it, try this (yes, really)
	command_string := "pulsemixer --list-sources | cut -f3 | grep ': UMC204HD' | cut -d ',' -f 1 | cut -c 6-"

	cmd := exec.Command("bash", "-c", command_string)
	out, err := cmd.Output()
	if err != nil {
		log.Warn().Msg("error")
		log.Warn().Err(err)
	}
	source = strings.TrimSpace(string(out))
	log.Info().Msg(source)
}

func (s *Mute) Buttons() {
	// Command
	shotbutton, _ := buttons.NewImageFileButton(viper.GetString("buttons.images") + "/mic.png")
	shotaction := &sdactionhandlers.CustomAction{}
	shotaction.SetHandler(func(btn streamdeck.Button) {
		go s.toggleMute()
	})
	shotbutton.SetActionHandler(shotaction)
	s.SD.AddButton(s.Button_id, shotbutton)
	s.updateButtonDecoration()
}

func (s *Mute) toggleMute() {
	if s.Status {
		// unmute
		log.Debug().Msg("Unmuting")
		cmd := exec.Command("pulsemixer", "--id", source, "--unmute")
		if err := cmd.Run(); err != nil {
			log.Warn().Err(err)
		}
	} else {
		log.Debug().Msg("Muting")
		cmd := exec.Command("pulsemixer", "--id", source, "--mute")
		if err := cmd.Run(); err != nil {
			log.Warn().Err(err)
		}
	}
	s.updateButtonDecoration()
}

func (s *Mute) readMuteStatus() bool {
	cmd := exec.Command("pulsemixer", "--id", source, "--get-mute")
	var outb bytes.Buffer
	cmd.Stdout = &outb
	if err := cmd.Run(); err != nil {
		log.Warn().Err(err)
	} else {
		// there's a newline in the stdout output!
		if outb.String() == "0\n" {
			log.Info().Msg("Mic is LIVE")
			s.Status = false
		} else {
			log.Info().Msg("Mic is muted")
			s.Status = true
		}
	}
	return s.Status

}

func (s *Mute) updateButtonDecoration() {
	status := s.readMuteStatus()
	decorate_on := sddecorators.NewBorder(12, color.RGBA{255, 120, 150, 255})
	decorate_off := sddecorators.NewBorder(12, color.RGBA{120, 120, 120, 255})
	if status {
		s.SD.SetDecorator(s.Button_id, decorate_off)
	} else {
		s.SD.SetDecorator(s.Button_id, decorate_on)
	}
}
