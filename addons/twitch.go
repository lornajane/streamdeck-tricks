package addons

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/magicmonkey/go-streamdeck"
	buttons "github.com/magicmonkey/go-streamdeck/buttons"
	"github.com/nicklaw5/helix"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Twitch struct {
	SD            *streamdeck.StreamDeck
	twitch_client helix.Client
}

func (t *Twitch) Init() {
	// make the twitch client
	client, err := helix.NewClient(&helix.Options{
		ClientID:     viper.GetString("twitch.client_id"),
		ClientSecret: viper.GetString("twitch.client_secret"),
		RedirectURI:  "http://localhost:7001/auth-callback",
	})
	if err != nil {
		panic(err)
	}
	t.twitch_client = *client

	// refresh token valid? Check by trying to use it to get new tokens
	isValid := t.updateTokens()

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

		resp, err := t.twitch_client.RequestUserAccessToken(code)
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

// UpdateTokens tries to use a current refresh token to get a new access token
// If we got new tokens, it returns true. If not (invalid or missing refresh token usually)
// then it returns false.
func (t *Twitch) updateTokens() bool {
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
	log.Debug().Msg("Twitch tokens updated")

	return true
}

func (t *Twitch) Buttons() {
	markbutton := buttons.NewTextButton("Mark")
	markbutton.SetActionHandler(&TwitchAction{Action: "mark", Client: t.twitch_client, Twitch: t})
	t.SD.AddButton(23, markbutton)
	vidbutton := buttons.NewTextButton("Vids")
	vidbutton.SetActionHandler(&TwitchAction{Action: "videos", Client: t.twitch_client, Twitch: t})
	t.SD.AddButton(22, vidbutton)
}

type TwitchAction struct {
	Client helix.Client
	Action string
	Twitch *Twitch
}

func (action *TwitchAction) Pressed(btn streamdeck.Button) {
	log.Info().Msg("Twitch Action: " + action.Action)

	log.Debug().Msg("Check access token (next line shows if we had to update it)")
	isValid, _, _ := action.Client.ValidateToken(action.Client.GetUserAccessToken())
	if !isValid {
		action.Twitch.updateTokens()
	}

	user_id := viper.GetString("twitch.user_id")

	if action.Action == "videos" {
		// output all videos for my user ID
		resp, err := action.Client.GetVideos(&helix.VideosParams{UserID: user_id})
		if err != nil {
			log.Error().Err(err)
		}
		fmt.Printf("%#v\n", resp)
	} else if action.Action == "mark" {
		// not going to do anything with these responses while I'm streaming
		resp_mark, _ := action.Client.CreateStreamMarker(&helix.CreateStreamMarkerParams{
			UserID:      user_id,
			Description: "Streamdeck marks the spot",
		})
		fmt.Printf("%#v\n", resp_mark)
	}
}
