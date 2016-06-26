package teamgen

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func isTokenValid(ctx context.Context, token string) bool {
	isValid := verificationToken == "" || token == verificationToken
	if !isValid {
		log.Errorf(ctx, "Invalid Slack token: %s", token)
	}
	return isValid
}

// handleCommand adds a member to a team
func handleCommand(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if !isTokenValid(ctx, r.PostFormValue("token")) {
		http.Error(w, "Invalid Slack token.", http.StatusBadRequest)
		return
	}

	cmdType, cmdArgs := parseCommand(r.PostFormValue("text"))
	resp := processComamnd(ctx, cmdType, cmdArgs, r.PostFormValue("team_id"), r.PostFormValue("channel_id"))

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Errorf(ctx, "Error encoding JSON: %s", err)
		http.Error(w, "Error encoding JSON.", http.StatusInternalServerError)
		return
	}
}

func handleSendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	// teamId := r.FormValue("teamId")
	channelID := r.FormValue("channelId")
	token := r.FormValue("token")

	// TODO: Get Token using teamID from OAuthAccessToken data
	data := url.Values{}
	data.Set("token", token)
	data.Set("channel", channelID)
	data.Set("text", "Hello world")
	encodedData := data.Encode()

	req, _ := http.NewRequest("POST", postMessageURL, bytes.NewBufferString(encodedData))
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var oauthResponse = new(OAuthAccessToken)
	err = json.Unmarshal(responseBody, &oauthResponse)
	if err != nil {
		panic(err.Error())
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(oauthResponse); err != nil {
		http.Error(w, "Error encoding JSON.", http.StatusInternalServerError)
		return
	}
}

func handleOauth(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Add("client_secret", clientSecret)
	data.Add("code", r.FormValue("code"))
	encodedData := data.Encode()

	req, _ := http.NewRequest("POST", oauthAccessURL, bytes.NewBufferString(encodedData))
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var oauthResponse = new(OAuthAccessToken)
	err = json.Unmarshal(responseBody, &oauthResponse)
	if err != nil {
		panic(err.Error())
	}

	key := generateOAuthAccessTokenKey(ctx, oauthResponse.TeamID)
	oauthResponse.LastUpdated = time.Now()
	if _, err := datastore.Put(ctx, key, oauthResponse); err != nil {
		log.Errorf(ctx, "Error on adding access token to datastore: %s", err)
		http.Error(w, "Error occurred during authorization process JSON.", http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(oauthResponse); err != nil {
		http.Error(w, "Error encoding JSON.", http.StatusInternalServerError)
		return
	}
}
