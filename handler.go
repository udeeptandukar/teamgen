package teamgen

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
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
	teamID := r.FormValue("teamId")
	channelID := r.FormValue("channelId")
	err := postMessage(ctx, teamID, channelID)
	if err != nil {
		log.Errorf(ctx, "Error on sending message: %s", err)
		fmt.Fprintf(w, "Could not send message: %s", err.Error())
	} else {
		fmt.Fprintf(w, "Message sent successfully to team %s, channel %s", teamID, channelID)
	}
}

func handleOauth(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	err := doOauthAuthorization(ctx, r.FormValue("code"))
	if err != nil {
		log.Errorf(ctx, "Error on authorization: %s", err)
		fmt.Fprintf(w, "Could not authorize: %s", err.Error())
	} else {
		fmt.Fprintf(w, "Team authorized successfully")
	}
}

func handleScheduling(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	teams, err := getAllTeams(ctx)
	if err != nil {
		log.Errorf(ctx, "Error on scheduling: %s", err)
		fmt.Fprintf(w, "Could not authorize: %s", err.Error())
	}

	for i := 0; i < len(teams); i++ {
		eta := time.Now().Add(time.Duration(8) * time.Hour)
		data := map[string][]string{"teamId": {teams[i].SlackTeamID}, "channelID": {teams[i].SlackChannelID}}
		t := taskqueue.NewPOSTTask("/sendMsg", data)
		t.ETA = eta
		if _, err := taskqueue.Add(ctx, t, "send-message"); err != nil {
			log.Errorf(ctx, "Error on scheduling: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
