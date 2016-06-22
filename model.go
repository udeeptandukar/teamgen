package teamgen

import (
	"time"

	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"
)

// TeamGenerator to store information required to generate random team
type TeamGenerator struct {
	Members        []string
	Schedules      []string
	NumberOfTeams  int
	RandomName     bool
	SlackTeamID    string
	SlackChannelID string
	LastUpdated    time.Time
}

// ChannelOAuth to store oauth token for a slack channel
type ChannelOAuth struct {
	SlackChannelID string
	Token          string
}

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

func generateTeamGeneratorKey(ctx context.Context, teamID string, channelID string) *datastore.Key {
	keyString := teamID + "-" + channelID
	return datastore.NewKey(ctx, "TeamGenerator", keyString, 0, nil)
}
