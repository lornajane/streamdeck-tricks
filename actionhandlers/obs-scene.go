package actionhandlers

import (
	"github.com/christopher-dG/go-obs-websocket"
	streamdeck "github.com/magicmonkey/go-streamdeck"
	"github.com/rs/zerolog/log"
)

type OBSSceneAction struct {
	Client obsws.Client
	Scene  string
	btn    streamdeck.Button
}

func (action *OBSSceneAction) Pressed(btn streamdeck.Button) {

	log.Info().Msg("Set scene: " + action.Scene)
	req := obsws.NewSetCurrentSceneRequest(action.Scene)
	_, err := req.SendReceive(action.Client)
	if err != nil {
		log.Warn().Err(err).Msg("OBS scene action error")
	}

}
