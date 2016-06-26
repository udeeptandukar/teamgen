package teamgen

import (
	"strings"

	"golang.org/x/net/context"
)

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
	default:
		resp = constructSlackCmdResponse("ephemeral", "Invalid command")
	}
	return resp
}
