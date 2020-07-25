package addons

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/magicmonkey/go-streamdeck"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Twitch struct {
	SD          *streamdeck.StreamDeck
	AccessToken string
}

// {"access_token":"wbffcg9rwe4p25sd6l2arp0u5rkzz5","expires_in":15680,"refresh_token":"wff8nn723olhosx2aqcl2613arum53gdon1k8rercvrq5281bx","scope":["user:edit:broadcast"],"token_type":"bearer"}

type TwitchRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (t *Twitch) Init() {
	t.refreshToken()

	/* --this bit doesn't work yet

	endpoint := "https://api.twitch.tv/helix/users?login=lornajanetv"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", endpoint, nil)
	fmt.Println(endpoint)
	fmt.Println(t.AccessToken)
	req.Header.Set("Authorization", "Bearer "+t.AccessToken)
	response, err := client.Do(req)
	if err != nil {
		log.Error().Err(err)
	}

	body, _ := ioutil.ReadAll(response.Body)
	log.Info().Msg(string(body))
	*/

}

// refreshToken takes the refresh token from config, gets a new access token and refresh token. The access token is set on the Twitch struct and the refresh token is written to config for next time.
func (t *Twitch) refreshToken() {
	values := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {viper.Get("twitch.refresh_token").(string)},
		"client_id":     {viper.Get("twitch.key").(string)},
		"client_secret": {viper.Get("twitch.secret").(string)},
	}

	response, err := http.PostForm("https://id.twitch.tv/oauth2/token", values)
	// response, err := http.PostForm("https://ljnexmo.eu.ngrok.io/", values)
	if err != nil {
		log.Error().Err(err)
	}

	body, _ := ioutil.ReadAll(response.Body)
	log.Info().Msg(string(body))

	// decode
	data := TwitchRefreshResponse{}
	data_err := json.Unmarshal(body, &data)
	if err != nil {
		log.Error().Err(data_err)
	} else {
		// store the access token for use in this program
		t.AccessToken = data.AccessToken
		// write the refresh token to the config file for another day or later
		viper.Set("twitch.refresh_token", data.RefreshToken)
		viper.WriteConfig()
	}
}
