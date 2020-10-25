package addons

import (
	"image/color"
	"strings"

	obsws "github.com/christopher-dG/go-obs-websocket"
	"github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/buttons"
	sddecorators "github.com/magicmonkey/go-streamdeck/decorators"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Obs struct {
	SD         *streamdeck.StreamDeck
	obs_client obsws.Client
	Offset     int
}

var obs_current_scene string

type ObsScene struct {
	Name     string `mapstructure:"name"`
	Image    string `mapstructure:"image"`
	ButtonId int
}

func (scene *ObsScene) SetButtonId(id int) {
	scene.ButtonId = id
}

var buttons_obs map[string]*ObsScene // scene name and image name

func (o *Obs) Init() {
	o.ConnectOBS()
	o.ObsEventHandlers()
}

func (o *Obs) ConnectOBS() {
	log.Debug().Msg("Connecting to OBS...")
	log.Info().Msgf("%#v\n", viper.Get("obs.host"))
	o.obs_client = obsws.Client{Host: "localhost", Port: 4444}
	err := o.obs_client.Connect()
	if err != nil {
		log.Warn().Err(err).Msg("Cannot connect to OBS")
	}
}

func (o *Obs) ObsEventHandlers() {
	if o.obs_client.Connected() == true {
		// Scene change
		o.obs_client.AddEventHandler("SwitchScenes", func(e obsws.Event) {
			// Make sure to assert the actual event type.
			scene := strings.ToLower(e.(obsws.SwitchScenesEvent).SceneName)
			log.Info().Msg("Old scene: " + obs_current_scene)
			// undecorate the old
			if oldb, ok := buttons_obs[obs_current_scene]; ok {
				log.Info().Int("button", oldb.ButtonId).Msg("Clear original button decoration")
				o.SD.UnsetDecorator(oldb.ButtonId)
			}
			// decorate the new
			log.Info().Msg("New scene: " + scene)
			if eventb, ok := buttons_obs[scene]; ok {
				decorator2 := sddecorators.NewBorder(5, color.RGBA{255, 0, 0, 255})
				log.Info().Int("button", eventb.ButtonId).Msg("Highlight new scene button")
				o.SD.SetDecorator(eventb.ButtonId, decorator2)
			}
			obs_current_scene = scene
		})

		// OBS Exits
		o.obs_client.AddEventHandler("Exiting", func(e obsws.Event) {
			log.Info().Msg("OBS has exited")
			o.ClearButtons()
		})

		// Scene Collection Switched
		o.obs_client.AddEventHandler("SceneCollectionChanged", func(e obsws.Event) {
			log.Info().Msg("Scene collection changed")
			o.ClearButtons()
			o.Buttons()
		})

	}
}

func (o *Obs) Buttons() {
	if o.obs_client.Connected() == true {
		// OBS Scenes to Buttons
		buttons_obs = make(map[string]*ObsScene)
		viper.UnmarshalKey("obs_scenes", &buttons_obs)
		image_path := viper.GetString("buttons.images")
		var image string

		// what scenes do we have? (max 8 for the top row of buttons)
		scene_req := obsws.NewGetSceneListRequest()
		scenes, err := scene_req.SendReceive(o.obs_client)
		if err != nil {
			log.Warn().Err(err)
		}
		obs_current_scene = strings.ToLower(scenes.CurrentScene)

		// make buttons for these scenes
		for i, scene := range scenes.Scenes {
			log.Debug().Msg("Scene: " + scene.Name)
			image = ""
			oaction := &OBSSceneAction{Scene: scene.Name, Client: o.obs_client}
			sceneName := strings.ToLower(scene.Name)

			if s, ok := buttons_obs[sceneName]; ok {
				if s.Image != "" {
					image = image_path + s.Image
				}
			} else {
				// there wasn't an entry in the buttons for this scene so add one
				buttons_obs[sceneName] = &ObsScene{}
			}

			if image != "" {
				// try to make an image button

				obutton, err := buttons.NewImageFileButton(image)
				if err == nil {
					obutton.SetActionHandler(oaction)
					o.SD.AddButton(i+o.Offset, obutton)
					// store which button we just set
					buttons_obs[sceneName].SetButtonId(i + o.Offset)
				} else {
					// something went wrong with the image, use a default one
					image = image_path + "/play.jpg"
					obutton, err := buttons.NewImageFileButton(image)
					if err == nil {
						obutton.SetActionHandler(oaction)
						o.SD.AddButton(i+o.Offset, obutton)
						// store which button we just set
						buttons_obs[sceneName].SetButtonId(i + o.Offset)
					}
				}
			} else {
				// use a text button
				oopbutton := buttons.NewTextButton(scene.Name)
				oopbutton.SetActionHandler(oaction)
				o.SD.AddButton(i+o.Offset, oopbutton)
				// store which button we just set
				buttons_obs[sceneName].SetButtonId(i + o.Offset)
			}

			// only need a few scenes
			if i > 5 {
				break
			}
		}

		// highlight the active scene
		if eventb, ok := buttons_obs[obs_current_scene]; ok {
			decorator2 := sddecorators.NewBorder(5, color.RGBA{255, 0, 0, 255})
			log.Info().Int("button", eventb.ButtonId).Msg("Highlight current scene")
			o.SD.SetDecorator(eventb.ButtonId, decorator2)
		}
	}

	// show a button to reinitialise all the OBS things
	startbutton := buttons.NewTextButton("Go OBS")
	startbutton.SetActionHandler(&OBSStartAction{Client: o.obs_client, Obs: o})
	o.SD.AddButton(o.Offset+7, startbutton)
}

func (o *Obs) ClearButtons() {
	for i := 0; i < 7; i++ {
		o.SD.UnsetDecorator(o.Offset + i)
		clearbutton := buttons.NewTextButton("")
		o.SD.AddButton(o.Offset+i, clearbutton)
	}
}

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

type OBSStartAction struct {
	Client obsws.Client
	Obs    *Obs
	btn    streamdeck.Button
}

func (action *OBSStartAction) Pressed(btn streamdeck.Button) {
	log.Debug().Msg("Reinit OBS")
	if !action.Obs.obs_client.Connected() {
		action.Obs.Init()
	}
	action.Obs.ClearButtons()
	action.Obs.Buttons()
}
