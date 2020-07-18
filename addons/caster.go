package addons

import (
	"image/color"
	"regexp"
	"strconv"

	"github.com/hypebeast/go-osc/osc"
	"github.com/lornajane/streamdeck-tricks/actionhandlers"
	"github.com/magicmonkey/go-streamdeck"
	buttons "github.com/magicmonkey/go-streamdeck/buttons"
	sddecorators "github.com/magicmonkey/go-streamdeck/decorators"
	"github.com/rs/zerolog/log"
)

type Caster struct {
	SD *streamdeck.StreamDeck
}

var buttons_osc map[int]string // just the track ID and name

func (c *Caster) Init() {
	go c.osc_server()
}

func (c *Caster) Buttons() {
	buttons_osc = make(map[int]string)
	osc_send_sync()
}

func (c *Caster) osc_server() {
	addr := "127.0.0.1:9000"
	d := osc.NewStandardDispatcher()
	d.AddMsgHandler("*", func(msg *osc.Message) {
		go c.osc_event(msg.Address, msg.Arguments)

	})

	server := &osc.Server{
		Addr:       addr,
		Dispatcher: d,
	}
	log.Debug().Msg("Starting OSC Listener")
	go server.ListenAndServe()
}

func (c *Caster) osc_event(Address string, Arguments []interface{}) {
	// buttons offset, where to start
	offset := 23

	// react to a new track name
	re := regexp.MustCompile(`/cbp/(.)/m/label/tr_name$`)
	info := re.FindStringSubmatch(Address)
	if len(info) > 1 {
		track_name := Arguments[0].(string)
		track_index, _ := strconv.Atoi(info[1])
		if (track_name) != "<Drop File>" {
			log.Debug().Msg("Track " + info[1] + " is: " + track_name)
			// make a button
			audiobutton := buttons.NewTextButton(track_name)
			audiobutton.SetActionHandler(&actionhandlers.OSCAction{Track: track_index})
			c.SD.AddButton(offset+track_index, audiobutton)
			buttons_osc[track_index] = track_name
		}
	}

	rf := regexp.MustCompile(`/cbp/(.)/m/label/p_s$`)
	info = rf.FindStringSubmatch(Address)
	if len(info) > 1 {
		track_index, _ := strconv.Atoi(info[1])
		// ignore all the tracks that we didn't register
		if _, ok := buttons_osc[track_index]; ok {
			action := Arguments[0].(string)
			if action == "Playing" {
				log.Debug().Msg("Playing: " + buttons_osc[track_index])
				decorator2 := sddecorators.NewBorder(12, color.RGBA{255, 255, 150, 255})
				c.SD.SetDecorator(offset+track_index, decorator2)
			}

			if action == "Stopped" {
				log.Debug().Msg("Stopped: " + buttons_osc[track_index])
				c.SD.UnsetDecorator(offset + track_index)
			}

		}
	}
}

func osc_send_sync() {
	client := osc.NewClient("127.0.0.1", 5051)
	msg := osc.NewMessage("/glo/sync")
	msg.Append(int32(1))
	err := client.Send(msg)
	if err != nil {
		log.Error().Err(err)
	}
}
