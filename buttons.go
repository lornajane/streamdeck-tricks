package main

import (
	"fmt"
	"os/exec"

	"github.com/magicmonkey/go-streamdeck"
	"github.com/spf13/viper"
)

var button_images [32]string

func MyButtonPress(btnIndex int, sd *streamdeck.Device) {
	switch btnIndex {
	case 7:
		cmd := exec.Command("xeyes")
		cmd.Start()
	default:
		ToggleImageOnButton(sd, btnIndex, viper.GetString("images_buttons")+"/play.jpg")
	}
}

func ToggleImageOnButton(sd *streamdeck.Device, btnIndex int, image string) {
	if button_images[btnIndex] == "" {
		sd.WriteImageToButton(image, btnIndex)
		button_images[btnIndex] = image
	} else {
		sd.WriteColorToButton(0, 0, 0, btnIndex)
		button_images[btnIndex] = ""
	}
}
