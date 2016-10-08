package teamgen

import (
	"fmt"
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
	case "enable":
		resp = toggleAutoGeneration(ctx, teamID, channelID, true)
	case "disable":
		resp = toggleAutoGeneration(ctx, teamID, channelID, false)
	case "generate":
		deferSendMsg(ctx, teamID, channelID, 0)
		resp = constructSlackCmdResponse("ephemeral", "Team generation scheduled.")
	case "help":
		helpMsg := fmt.Sprintf("Use `/%s` to generate random teams of pairs\n", slackCommand)
		helpMsg += fmt.Sprintf("• `/%s member-add Joe Iris Doe` to add members\n", slackCommand)
		helpMsg += fmt.Sprintf("• `/%s member-exclusion Joe,Iris Joe,Doe` to have members not in same team\n", slackCommand)
		helpMsg += fmt.Sprintf("• `/%s show-config` to show current configurations\n", slackCommand)
		helpMsg += fmt.Sprintf("• `/%s generate` to generate teams\n", slackCommand)
		helpMsg += fmt.Sprintf("For more about %s, visit <https://pair-generator.appspot.com/|here>.", applicationName)
		resp = constructSlackCmdResponse("ephemeral", helpMsg)
	default:
		resp = constructSlackCmdResponse("ephemeral", "Invalid command")
	}
	return resp
}
