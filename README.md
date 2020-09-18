# StreamDeck Tricks

This repo holds the one-of-a-kind project that drives my streamdeck application. It is not intended for general use, and is just released to allow community members to improve their own knowledge and tooling. That said, I'm happy to review any issues and pull requests, but please be aware that I'm unlikely to merge or implement any changes that don't help my own workflow.

If you use this repository to improve your streamdeck workflow, you are encouraged to do the same thing - share your configuration and let others know how you made your stream more awesome!

## Fabulous Libraries That Made This Project Possible

* [go-streamdeck](https://github.com/magicmonkey/go-streamdeck) is the key library to drive the streamdeck part of the application.
* [hypebeast/osc](https://github.com/hypebeast/go-osc) Go library for using Open Sound Control applications - I use mine with [CasterSoundboard](https://github.com/JupiterBroadcasting/CasterSoundboard/) on Ubuntu.
* [paho.mqtt.golang](https://github.com/eclipse/paho.mqtt.golang) for MQTT integration (with my [LED shelf](https://lornajane.net/posts/2020/neopixel-shelf).
* [OBS Websockets](https://github.com/christopher-dG/go-obs-websocket) to change scenes and get current state information about the selected scene.
* [Helix](https://github.com/nicklaw5/helix) for Twitch API integration. It doesn't do much with chat but is useful for metadata and I drop markers while I'm streaming.
* Not really a library but I do also use the buttons to call individual commands. You can see an example in my [blog post about a screenshot button](https://lornajane.net/posts/2020/add-a-screenshot-button-to-streamdeck-with-golang)
