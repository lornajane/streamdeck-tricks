package main

import (
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/lornajane/streamdeck-tricks/addons"
	streamdeck "github.com/magicmonkey/go-streamdeck"
	_ "github.com/magicmonkey/go-streamdeck/devices"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var sd *streamdeck.StreamDeck

func loadConfigAndDefaults() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04"})

	// first set some default values
	viper.AddConfigPath(".")
	viper.SetDefault("buttons.images", "images/buttons") // location of button images
	viper.SetDefault("obs.host", "localhost")            // OBS webhooks endpoint
	viper.SetDefault("obs.port", 4444)                   // OBS webhooks endpoint
	viper.SetDefault("mqtt.uri", "tcp://10.1.0.1:1883")  // MQTT server location

	// now read in config for any overrides
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Warn().Msgf("Cannot read config file: %s \n", err)
	}

	// useful in development phase, pick up config file updates
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info().Msgf("Config file changed:", e.Name)
	})

}

func main() {
	loadConfigAndDefaults()
	log.Info().Msg("Starting streamdeck tricks. Hai!")

	var err error
	sd, err = streamdeck.New()
	if err != nil {
		log.Error().Err(err).Msg("Error finding Streamdeck")
		panic(err)
	}

	// init OBS
	// Initialise OBS to use OBS features (requires websockets plugin in OBS)
	obs_addon := addons.Obs{SD: sd}
	obs_addon.ConnectOBS()
	obs_addon.ObsEventHandlers()
	obs_addon.Buttons()

	// init MQTT
	mqtt_addon := addons.MqttThing{SD: sd}
	mqtt_addon.Init()
	mqtt_addon.Buttons()

	// init Screenshot
	screenshot_addon := addons.Screenshot{SD: sd}
	screenshot_addon.Init()
	screenshot_addon.Buttons()

	// set up soundcaster
	caster_addon := addons.Caster{SD: sd}
	caster_addon.Init()
	caster_addon.Buttons()

	log.Info().Msg("Up and running")
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
