package main

import (
	"image/color"
	"strings"

	obsws "github.com/christopher-dG/go-obs-websocket"
	"github.com/lornajane/streamdeck-tricks/actionhandlers"
	// "github.com/lornajane/streamdeck-tricks/addons"

	buttons "github.com/magicmonkey/go-streamdeck/buttons"
	sddecorators "github.com/magicmonkey/go-streamdeck/decorators"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	_ "github.com/godbus/dbus"
)

var obs_client obsws.Client
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

// InitButtons sets up initial button prompts
func InitButtons() {

	// Initialise OBS to use OBS features (requires websockets plugin in OBS)
	obs_client = connectOBS()

	if obs_client.Connected() == true {
		obs_client.AddEventHandler("SwitchScenes", func(e obsws.Event) {
			// Make sure to assert the actual event type.
			scene := strings.ToLower(e.(obsws.SwitchScenesEvent).SceneName)
			log.Info().Msg("Old scene: " + obs_current_scene)
			// undecorate the old
			if oldb, ok := buttons_obs[obs_current_scene]; ok {
				log.Info().Int("button", oldb.ButtonId).Msg("Clear original button decoration")
				sd.UnsetDecorator(oldb.ButtonId)
			}
			// decorate the new
			log.Info().Msg("New scene: " + scene)
			if eventb, ok := buttons_obs[scene]; ok {
				decorator2 := sddecorators.NewBorder(5, color.RGBA{255, 0, 0, 255})
				log.Info().Int("button", eventb.ButtonId).Msg("Highlight new scene button")
				sd.SetDecorator(eventb.ButtonId, decorator2)
			}
			obs_current_scene = scene
		})
	}

	// OBS Scenes to Buttons
	buttons_obs = make(map[string]*ObsScene)
	viper.UnmarshalKey("obs_scenes", &buttons_obs)

	if obs_client.Connected() == true {
		// offset for what number button to start at
		offset := 0
		image_path := viper.GetString("buttons.images")
		var image string

		// what scenes do we have? (max 8 for the top row of buttons)
		scene_req := obsws.NewGetSceneListRequest()
		scenes, err := scene_req.SendReceive(obs_client)
		if err != nil {
			log.Warn().Err(err)
		}
		obs_current_scene = strings.ToLower(scenes.CurrentScene)

		// make buttons for these scenes
		for i, scene := range scenes.Scenes {
			log.Debug().Msg("Scene: " + scene.Name)
			image = ""
			oaction := &actionhandlers.OBSSceneAction{Scene: scene.Name, Client: obs_client}
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
					sd.AddButton(i+offset, obutton)
					// store which button we just set
					buttons_obs[sceneName].SetButtonId(i + offset)
				} else {
					// something went wrong with the image, use a default one
					image = image_path + "/play.jpg"
					obutton, err := buttons.NewImageFileButton(image)
					if err == nil {
						obutton.SetActionHandler(oaction)
						sd.AddButton(i+offset, obutton)
						// store which button we just set
						buttons_obs[sceneName].SetButtonId(i + offset)
					}
				}
			} else {
				// use a text button
				oopbutton := buttons.NewTextButton(scene.Name)
				oopbutton.SetActionHandler(oaction)
				sd.AddButton(i+offset, oopbutton)
				// store which button we just set
				buttons_obs[sceneName].SetButtonId(i + offset)
			}

			// only need a few scenes
			if i > 6 {
				break
			}
		}

		// highlight the active scene
		if eventb, ok := buttons_obs[obs_current_scene]; ok {
			decorator2 := sddecorators.NewBorder(5, color.RGBA{255, 0, 0, 255})
			log.Info().Int("button", eventb.ButtonId).Msg("Highlight current scene")
			sd.SetDecorator(eventb.ButtonId, decorator2)
		}

	}

}

func connectOBS() obsws.Client {
	log.Debug().Msg("Connecting to OBS...")
	log.Info().Msgf("%#v\n", viper.Get("obs.host"))
	obs_client = obsws.Client{Host: "localhost", Port: 4444}
	err := obs_client.Connect()
	if err != nil {
		log.Warn().Err(err).Msg("Cannot connect to OBS")
	}
	return obs_client
}
