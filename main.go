package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/eiannone/keyboard"
	"github.com/lornajane/streamdeck-tricks/addons"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func loadConfigAndDefaults() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04"})

	// first set some default values
	viper.AddConfigPath(".")

	// now read in config for any overrides
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Warn().Msgf("Cannot read config file: %s \n", err)
	}
}

func captureKeyboard(mqtt_addon *addons.MqttThing) {
	if err := keyboard.Open(); err != nil {
		log.Error().Err(err).Msg("Failed to open keyboard")
		return
	}
	defer keyboard.Close()

	log.Info().Msg("Press keys 0-9 (ESC to quit)")
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Error().Err(err).Msg("Keyboard error")
			break
		}

		if key == keyboard.KeyEsc {
			log.Info().Msg("ESC pressed, exiting")
			break
		}

		if char >= '0' && char <= '9' {
			log.Info().Msgf("Key pressed: %c", char)
			mqtt_addon.KeyPress(char)
		}
	}
}

func main() {
	loadConfigAndDefaults()
	log.Info().Msg("Starting LEDWall tricks. Hai!")

	// init MQTT
	mqtt_addon := addons.MqttThing{}
	mqtt_addon.Init()

	go webserver()
	go captureKeyboard(&mqtt_addon)

	log.Info().Msg("Up and running")
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func webserver() {
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	http.ListenAndServe(":7001", nil)
}
