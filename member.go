package teamgen

import (
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func addMember(ctx context.Context, teamID string, channelID string, members []string) SlashResponse {
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
		return constructSlashResponse("ephemeral", "Team members added: "+strings.Join(members, ", "))
	}

	return constructSlashResponse("ephemeral", "Error occurred while adding members. Please try again.")
}

func showConfig(ctx context.Context, teamID string, channelID string) SlashResponse {
	key := generateTeamGeneratorKey(ctx, teamID, channelID)
	teamGenerator := new(TeamGenerator)

	if err := datastore.Get(ctx, key, teamGenerator); err == nil {
		text := "Team members: " + strings.Join(teamGenerator.Members, ", ")
		text += "\nNo. of teams: " + strconv.Itoa(teamGenerator.NumberOfTeams)
		return constructSlashResponse("ephemeral", text)
	}

	return constructSlashResponse("ephemeral", "No config found.")
}
