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
	NumberOfTeams  int16
	RandomName     bool
	SlackTeamID    string
	SlackChannelID string
	LastUpdated    time.Time
}

type slashResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func constructSlashResponse(responseType string, text string) slashResponse {
	return slashResponse{
		ResponseType: responseType,
		Text:         text,
	}
}

func generateTeamGeneratorKey(ctx context.Context, teamID string, channelID string) *datastore.Key {
	keyString := teamID + "-" + channelID
	return datastore.NewKey(ctx, "TeamGenerator", keyString, 0, nil)
}
