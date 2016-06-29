package teamgen

import (
	"strings"

	"golang.org/x/net/context"
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
	switch cmdType {
	case "member-add":
		resp = addMember(ctx, teamID, channelID, args)
	case "show-config":
		resp = showConfig(ctx, teamID, channelID)
	case "generate":
		if err := postMessage(ctx, teamID, channelID); err != nil {
			resp = constructSlackCmdResponse("ephemeral", "Something went wrong. Please try again.")
		} else {
			resp = constructSlackCmdResponse("ephemeral", "Team is generated.")
		}
	default:
		resp = constructSlackCmdResponse("ephemeral", "Invalid command")
	}
	return resp
}
