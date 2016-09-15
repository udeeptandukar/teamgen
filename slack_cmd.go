package teamgen

import (
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

// SlackCmdResponse is slash response
type SlackCmdResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func constructSlackCmdResponse(responseType string, text string) SlackCmdResponse {
	return SlackCmdResponse{
		ResponseType: responseType,
		Text:         text,
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

func processComamnd(ctx context.Context, cmdType string, args []string, teamID string, channelID string) SlackCmdResponse {
	var resp SlackCmdResponse
	log.Debugf(ctx, "Processing command: %s", cmdType)

	switch cmdType {
	case "member-add":
		resp = addMember(ctx, teamID, channelID, args)
	case "member-exclusion":
		resp = addMemberExclusions(ctx, teamID, channelID, args)
	case "show-config":
		resp = showConfig(ctx, teamID, channelID)
	case "generate":
		msg, err := postMessage(ctx, teamID, channelID)
		if err != nil {
			log.Errorf(ctx, "Error generating team: %s", err)
			resp = constructSlackCmdResponse("ephemeral", "Team not generated. Please try again.")
		} else {
			resp = constructSlackCmdResponse("ephemeral", msg)
		}
	case "help":
		helpMsg := "Use `/pair-generator` to generate random teams of pairs\n"
		helpMsg += "• `/pair-generator member-add Joe Iris Doe` to add members\n"
		helpMsg += "• `/pair-generator member-exclusion Joe Iris` to have members not in same team\n"
		helpMsg += "• `/pair-generator show-config` to show current configurations\n"
		helpMsg += "• `/pair-generator generate` to generate teams\n"
		helpMsg += "For more about Pair Generator, visit <https://pair-generator.appspot.com/|here>."
		resp = constructSlackCmdResponse("ephemeral", helpMsg)
	default:
		resp = constructSlackCmdResponse("ephemeral", "Invalid command")
	}
	return resp
}
