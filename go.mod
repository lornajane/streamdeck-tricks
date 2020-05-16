module github.com/lornajane/streamdeck

go 1.14

//  require github.com/magicmonkey/go-streamdeck v0.0.0-20200514153614-0a6d100a5cec // indirect
replace github.com/magicmonkey/go-streamdeck => /home/lorna/go/src/github.com/magicmonkey/go-streamdeck

require (
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/fsnotify/fsnotify v1.4.7
	github.com/magicmonkey/go-streamdeck v0.0.0-00010101000000-000000000000
	github.com/spf13/viper v1.7.0
	golang.org/x/net v0.0.0-20200513185701-a91f0712d120 // indirect
)
