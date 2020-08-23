package addons

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/magicmonkey/go-streamdeck"
	buttons "github.com/magicmonkey/go-streamdeck/buttons"
	"github.com/nicklaw5/helix"
	"github.com/rs/zerolog/log"
)

type Twitch struct {
	SD            *streamdeck.StreamDeck
	twitch_client helix.Client
}

func (t *Twitch) Init() {
	// make the twitch client
	client, err := helix.NewClient(&helix.Options{
		ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		RedirectURI:  "http://localhost:7001/auth-callback",
	})
	if err != nil {
		panic(err)
	}
	t.twitch_client = *client

	// refresh token valid?
	isValid := t.refreshToken()

	if !isValid {
		// refresh token outdated or missing, re-auth
		fmt.Println("Auth to Twitch with URL in browser:")
		// now set up the auth URL
		scopes := []string{"user:edit:broadcast"}
		url := t.twitch_client.GetAuthorizationURL(&helix.AuthorizationURLParams{
			ResponseType: "code", // or "token"
			Scopes:       scopes,
			ForceVerify:  false,
		})
		fmt.Printf("%s\n", url)
	}

	// add the HTTP endpoint for the auth callback
	http.HandleFunc("/auth-callback", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "authed, like truthed")

		code := r.URL.Query().Get("code")

		resp, err := t.twitch_client.GetUserAccessToken(code)
		if err != nil {
			panic(err)
		}

		access_token := resp.Data.AccessToken
		// Set the access token on the client
		client.SetUserAccessToken(access_token)

		refresh_token := resp.Data.RefreshToken
		// Put the refresh token in a file for later
		ioutil.WriteFile("twitch_refresh_token", []byte(refresh_token), 0644)
	})
}

func (t *Twitch) refreshToken() bool {
	data, file_err := ioutil.ReadFile("twitch_refresh_token")
	if file_err != nil {
		log.Error().Msg("Cannot read twitch_refresh_token, not authed")
		return false
	}
	refreshToken := string(data)

	resp, t_err := t.twitch_client.RefreshUserAccessToken(refreshToken)
	if t_err != nil {
		log.Error().Msg("Could not refesh tokens, not authed")
		return false
	}

	access_token := resp.Data.AccessToken
	refresh_token := resp.Data.RefreshToken
	ioutil.WriteFile("twitch_refresh_token", []byte(refresh_token), 0644)
	// Set the access token on the client
	t.twitch_client.SetUserAccessToken(access_token)
	log.Debug().Msg("Token refreshed")

	return true
}

func (t *Twitch) Buttons() {
	markbutton := buttons.NewTextButton("Mark")
	markbutton.SetActionHandler(&TwitchAction{Action: "mark", Client: t.twitch_client})
	t.SD.AddButton(23, markbutton)
	vidbutton := buttons.NewTextButton("Vids")
	vidbutton.SetActionHandler(&TwitchAction{Action: "videos", Client: t.twitch_client})
	t.SD.AddButton(22, vidbutton)
}

type TwitchAction struct {
	Client helix.Client
	Action string
}

func (action *TwitchAction) Pressed(btn streamdeck.Button) {
	log.Debug().Msg("Twitch Action: " + action.Action)

	// _ := t.refreshToken()

	if action.Action == "videos" {
		// output all videos for my user ID
		resp, err := action.Client.GetVideos(&helix.VideosParams{UserID: "493107973"})
		if err != nil {
			log.Error().Err(err)
		}
		fmt.Printf("%#v\n", resp)
	} else if action.Action == "mark" {
		// not going to do anything with these responses while I'm streaming
		resp_mark, _ := action.Client.CreateStreamMarker(&helix.CreateStreamMarkerParams{
			UserID:      "493107973",
			Description: "Streamdeck marks the spot",
		})
		fmt.Printf("%#v\n", resp_mark)
	}
}
