package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	streamdeck "github.com/magicmonkey/go-streamdeck"
	"github.com/spf13/viper"
)

func loadConfigAndDefaults() {
	viper.AddConfigPath(".")

	// first set some default values
	viper.SetDefault("images_buttons", "images/buttons") // location of button images

	// now read in config for any overrides
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		fmt.Printf("Cannot read config file: %s \n", err)
	}
	fmt.Println(viper.Get("images_buttons"))

	// useful in development phase, pick up config file updates
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
}

func main() {
	loadConfigAndDefaults()

	sd := streamdeck.Open()
	sd.ClearButtons()

	sd.SetBrightness(50)
	InitButtons(sd)
	sd.ButtonPress(MyButtonPress)
}
