package teamgen

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
)

type response struct {
	Ok    bool
	Error string
}

// Attachments holds structure to send attachments to slack
type Attachments struct {
	Color   string `json:"color"`
	PreText string `json:"pretext"`
	Text    string `json:"text"`
}

func postMessage(ctx context.Context, teamID string, channelID string) error {
	client := urlfetch.Client(ctx)

	token, err := getBotAccessToken(ctx, teamID)
	if err != nil {
		panic(err.Error())
	}

	data := url.Values{}
	data.Set("token", token)
	data.Set("channel", channelID)
	data.Set("text", "Randomly selected teams:")
	teams, err := getRandomTeams(ctx, teamID, channelID)
	if err != nil {
		return err
	}
	attachments := []Attachments{}
	colors := []string{"#004d99", "#00cc00", "#660066", "#ff3300"}
	for i := 0; i < len(teams); i++ {
		attachments = append(attachments, Attachments{
			Color: colors[i%len(colors)],
			Text:  teams[i],
		})
	}
	ats, err := json.Marshal(attachments)
	if err != nil {
		return err
	}
	data.Set("attachments", string(ats))
	encodedData := data.Encode()

	req, _ := http.NewRequest("POST", postMessageURL, bytes.NewBufferString(encodedData))
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	respBody := &response{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(respBody); err != nil {
		return err
	}

	if !respBody.Ok {
		return errors.New(respBody.Error)
	}

	return nil
}

func doOauthAuthorization(ctx context.Context, code string) error {
	client := urlfetch.Client(ctx)

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Add("client_secret", clientSecret)
	data.Add("code", code)
	encodedData := data.Encode()

	req, _ := http.NewRequest("POST", oauthAccessURL, bytes.NewBufferString(encodedData))
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var oauthResponse = new(OAuthAccessToken)
	if err := json.Unmarshal(responseBody, &oauthResponse); err != nil {
		return err
	}

	key := generateOAuthAccessTokenKey(ctx, oauthResponse.TeamID)
	oauthResponse.LastUpdated = time.Now()
	if _, err := datastore.Put(ctx, key, oauthResponse); err != nil {
		return err
	}

	return nil
}
