package addons

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/atotto/clipboard"
	"github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/buttons"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Nightbot struct {
	SD           *streamdeck.StreamDeck
	AccessToken  string
	RefreshToken string
	Text         string
	ChatLines    []string
	ChatIndex    int
}

type NightbotAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func (n *Nightbot) Init() {
	// can we use a refresh token to avoid needing to click?
	go n.updateTokens("refresh_token", "")

	// add the HTTP endpoint for the auth callback
	http.HandleFunc("/nightbot", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK, auth in progress")
		code := r.URL.Query().Get("code")
		log.Debug().Msg(code)

		// cash in the code for a token, a responsible adult would check the state/nonce too
		n.updateTokens("code", code)
	})

	n.readChatLines()
}

// updateTokens can take either a "code" or a "refresh_token" with no value (it's stored on disk) to get new tokens
func (n *Nightbot) updateTokens(token_type string, code string) bool {
	client_id := viper.GetString("nightbot.client_id")
	client_secret := viper.GetString("nightbot.client_secret")
	redirect_uri := "https://agfa87agne4.eu.ngrok.io/nightbot"

	token_url := "https://api.nightbot.tv/oauth2/token"
	values := url.Values{}
	values.Set("client_id", client_id)
	values.Set("client_secret", client_secret)
	values.Set("redirect_uri", redirect_uri)

	if token_type == "code" {
		values.Set("grant_type", "authorization_code")
		values.Set("code", code)
	} else if token_type == "refresh_token" {
		data, file_err := ioutil.ReadFile("nightbot_refresh_token")
		if file_err != nil {
			log.Error().Msg("Cannot read nightbot_refresh_token, not authed")
			return false
		}
		values.Set("refresh_token", string(data))
		values.Set("grant_type", "refresh_token")
	}

	resp, err := http.PostForm(token_url, values)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	/*
		// Useful debugging but leaks creds, don't use when streaming
		log.Debug().Msg(string(body))
	*/

	if resp.StatusCode == 200 {
		var tokens NightbotAuthTokenResponse
		body, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(body, &tokens)
		if err != nil {
			log.Error().Err(err)
		}

		// all good!
		log.Info().Msg("Nightbot token refreshed")

		n.AccessToken = tokens.AccessToken
		n.RefreshToken = tokens.RefreshToken

		// Put the refresh token in a file for later
		ioutil.WriteFile("nightbot_refresh_token", []byte(tokens.RefreshToken), 0644)
		resp.Body.Close()
		return true
	} else {
		// something went wrong, build the auth URL and show it
		rand.Seed(time.Now().UnixNano())
		nonce := strconv.Itoa(rand.Intn(999))
		auth_url := "https://api.nightbot.tv/oauth2/authorize" +
			"?client_id=" + client_id +
			"&redirect_uri=" + redirect_uri +
			"&scope=channel_send" +
			"&response_type=code" +
			"&state=" + nonce

		log.Info().Msg("Click to auth for nightbot")
		log.Info().Msg(auth_url)

	}

	resp.Body.Close()
	return false
}

func (n *Nightbot) Buttons() {
	cuebutton := buttons.NewTextButton("Cue")
	cuebutton.SetActionHandler(&NightbotAction{Action: "chat-cue", Bot: n})
	n.SD.AddButton(9, cuebutton)

	chatsendbutton := buttons.NewTextButton("Send")
	chatsendbutton.SetActionHandler(&NightbotAction{Action: "chat-send", Bot: n})
	n.SD.AddButton(8, chatsendbutton)

	filebutton := buttons.NewTextButton("File")
	filebutton.SetActionHandler(&NightbotAction{Action: "chat-file", Bot: n})
	n.SD.AddButton(10, filebutton)
}

func (n *Nightbot) readChatLines() {
	var chatlines []string
	file, err := os.Open("motd.txt")
	if err != nil {
		log.Error().Err(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		chatlines = append(chatlines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Error().Err(err)
	}
	n.ChatLines = chatlines
	n.ChatIndex = 0
}

type NightbotAction struct {
	Action  string
	Bot     *Nightbot
	Scanner bufio.Scanner
}

type NightbotMessage struct {
	Message string `json:"message"`
}

func (action *NightbotAction) Pressed(btn streamdeck.Button) {
	log.Info().Msg("Bot Action: " + action.Action)

	if action.Action == "chat-cue" {
		text, _ := clipboard.ReadAll()
		action.Bot.Text = text
		log.Info().Msg("Ready to Send: " + action.Bot.Text)
	}

	if action.Action == "chat-send" {
		log.Info().Msg("Sending: " + action.Bot.Text)
		msg := NightbotMessage{Message: action.Bot.Text}
		json_data, _ := json.Marshal(msg)
		req_data := bytes.NewReader(json_data)

		chat_url := "https://api.nightbot.tv/1/channel/send"
		req, err := http.NewRequest("POST", chat_url, req_data)
		req.Header.Add("Authorization", "Bearer "+action.Bot.AccessToken)
		req.Header.Add("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		log.Info().Msg(string(body))
	}

	if action.Action == "chat-file" {
		text := action.Bot.ChatLines[action.Bot.ChatIndex]
		log.Debug().Msg("Ready to send: " + text)
		action.Bot.Text = text
		action.Bot.ChatIndex++
		if action.Bot.ChatIndex >= len(action.Bot.ChatLines) {
			action.Bot.readChatLines()
		}
	}

}
