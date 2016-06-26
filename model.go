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

func generateTeamGeneratorKey(ctx context.Context, teamID string, channelID string) *datastore.Key {
	keyString := teamID + "-" + channelID
	return datastore.NewKey(ctx, "TeamGenerator", keyString, 0, nil)
}

// BotUser holds structure for bot user and access token
type BotUser struct {
	BotUserID      string `json:"bot_user_id"`
	BotAccessToken string `json:"bot_access_token"`
}

// OAuthAccessToken to store oauth token for a slack channel
type OAuthAccessToken struct {
	AccessToken string  `json:"access_token"`
	Scope       string  `json:"scope"`
	TeamID      string  `json:"team_id"`
	TeamName    string  `json:"team_name"`
	Bot         BotUser `json:"bot"`
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

func generateOAuthAccessTokenKey(ctx context.Context, teamID string) *datastore.Key {
	return datastore.NewKey(ctx, "OAuthAccessToken", teamID, 0, nil)
}
