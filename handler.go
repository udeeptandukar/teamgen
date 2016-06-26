package teamgen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
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
	// ctx := appengine.NewContext(r)
	// teamId := r.FormValue("teamId")
	// channelID := r.FormValue("channelId")
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

	fmt.Fprint(w, resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var bodydata interface{}
	err = json.Unmarshal(body, &bodydata)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprint(w, bodydata)

	// w.Header().Set("content-type", "application/json")
	// if err := json.NewEncoder(w).Encode(resp.Body); err != nil {
	// 	log.Errorf(ctx, "Error encoding JSON: %s", err)
	// 	http.Error(w, "Error encoding JSON.", http.StatusInternalServerError)
	// 	return
	// }
}
