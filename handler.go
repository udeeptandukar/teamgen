package teamgen

import (
	"encoding/json"
	"net/http"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// handleCommand adds a member to a team
func handleCommand(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	postedToken := r.PostFormValue("token")
	if token != "" && postedToken != token {
		log.Errorf(ctx, "Invalid Slack token: %s", postedToken)
		http.Error(w, "Invalid Slack token.", http.StatusBadRequest)
		return
	}

	cmdType, cmdArgs := parseCommand(r.PostFormValue("text"))
	resp := processComamnd(cmdType, cmdArgs, r.PostFormValue("team_id"), r.PostFormValue("channel_id"))

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Errorf(ctx, "Error encoding JSON: %s", err)
		http.Error(w, "Error encoding JSON.", http.StatusInternalServerError)
		return
	}
}

func parseCommand(command string) (string, []string) {
	result := strings.Fields(command)
	cmdType := result[0]
	cmdArgs := []string{}
	if len(result) > 1 {
		cmdArgs = result[1:]
	}
	return cmdType, cmdArgs
}

func processComamnd(cmdType string, args []string, teamID string, channelID string) slashResponse {
	resp := constructSlashResponse("ephemeral", "Invalid command")
	if cmdType == "member-add" {
		resp = addMember(teamID, channelID, args)
	}
	return resp
}
