package teamgen

import (
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func addMember(ctx context.Context, teamID string, channelID string, members []string) SlackCmdResponse {
	key := generateTeamGeneratorKey(ctx, teamID, channelID)
	teamGenerator := new(TeamGenerator)

	if err := datastore.Get(ctx, key, teamGenerator); err != nil {
		teamGenerator.SlackTeamID = teamID
		teamGenerator.SlackChannelID = channelID
		teamGenerator.Members = members
		teamGenerator.RandomName = false
		teamGenerator.NumberOfTeams = len(members) / 2
		teamGenerator.LastUpdated = time.Now()
	} else {
		teamGenerator.Members = members
		teamGenerator.LastUpdated = time.Now()
	}

	if _, err := datastore.Put(ctx, key, teamGenerator); err == nil {
		return constructSlackCmdResponse("ephemeral", "Team members added: "+strings.Join(members, ", "))
	}

	return constructSlackCmdResponse("ephemeral", "Error occurred while adding members. Please try again.")
}

func showConfig(ctx context.Context, teamID string, channelID string) SlackCmdResponse {
	key := generateTeamGeneratorKey(ctx, teamID, channelID)
	teamGenerator := new(TeamGenerator)

	if err := datastore.Get(ctx, key, teamGenerator); err == nil {
		text := "Team members: " + strings.Join(teamGenerator.Members, ", ")
		text += "\nNo. of teams: " + strconv.Itoa(teamGenerator.NumberOfTeams)
		return constructSlackCmdResponse("ephemeral", text)
	}

	return constructSlackCmdResponse("ephemeral", "No config found.")
}
