package addons

import (
//	"image/color"
	"strings"

	"github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/buttons"
	// sddecorators "github.com/magicmonkey/go-streamdeck/decorators"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/scenes"
)

type Obs struct {
	SD         *streamdeck.StreamDeck
	obs_client *goobs.Client
	Offset     int
	connected  bool
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
	o.connected = false
	o.ConnectOBS()
}

func (o *Obs) ConnectOBS() {
	var err error
	log.Debug().Msg("Connecting to OBS...")
	
	o.obs_client, err = goobs.New("localhost:4455")
	if err != nil {
		log.Warn().Err(err).Msg("Cannot connect to OBS")
	} else {
		o.connected = true
	}
}

func (o *Obs) Buttons() {
	if o.connected {
		// OBS Scenes to Buttons
		buttons_obs = make(map[string]*ObsScene)
		viper.UnmarshalKey("obs_scenes", &buttons_obs)
		image_path := viper.GetString("buttons.images")
		var image string

		// what scenes do we have? (max 8 for the top row of buttons)
		scenes, _:= o.obs_client.Scenes.GetSceneList()
		for _, v := range scenes.Scenes {
			log.Debug().Msg("Found scene "+ v.SceneName)
		}
		obs_current_scene = strings.ToLower(scenes.CurrentProgramSceneName)
		log.Debug().Msg("Current scene: " + obs_current_scene)

		// make buttons for these scenes
		for i, scene := range scenes.Scenes {
			log.Debug().Msg("Scene: " + scene.SceneName)
			image = ""
			oaction := &OBSSceneAction{Scene: scene.SceneName, Client: o.obs_client}
			sceneName := strings.ToLower(scene.SceneName)

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
				oopbutton := buttons.NewTextButton(scene.SceneName)
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

	}
	// show a button to reinitialise all the OBS things

}

type OBSSceneAction struct {
	Client *goobs.Client
	Scene  string
	btn    streamdeck.Button
}

func (action *OBSSceneAction) Pressed(btn streamdeck.Button) {

	log.Info().Msg("Set scene: " + action.Scene)
	_, err := action.Client.Scenes.SetCurrentProgramScene(
		&scenes.SetCurrentProgramSceneParams{SceneName: action.Scene})
	if err != nil {
		log.Debug().Msg("oh dear")
	}
}


