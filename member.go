package teamgen

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func addMember(ctx context.Context, teamID string, channelID string, members []string) slashResponse {
	key := generateTeamGeneratorKey(ctx, teamID, channelID)

	// var teamGenerator TeamGenerator
	// if err := datastore.Get(ctx, key, &teamGenerator); err != nil {
	// 	teamGenerator := &TeamGenerator{
	// 		SlackTeamID:    teamID,
	// 		SlackChannelID: channelID,
	// 		Members:        members,
	// 		LastUpdated:    time.Now(),
	// 	}
	// } else {
	// 	teamGenerator.Members = members
	// 	teamGenerator.LastUpdated = time.Now()
	// }
	teamGenerator := &TeamGenerator{
		SlackTeamID:    teamID,
		SlackChannelID: channelID,
		Members:        members,
		LastUpdated:    time.Now(),
	}

	if _, err := datastore.Put(ctx, key, teamGenerator); err == nil {
		return constructSlashResponse("ephemeral", "Team members added")
	}

	return constructSlashResponse("ephemeral", "Error occurred while adding members. Please try again.")
}
