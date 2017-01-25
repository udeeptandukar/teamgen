package teamgen

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

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
	LastUpdated time.Time
}

func generateOAuthAccessTokenKey(ctx context.Context, teamID string) *datastore.Key {
	return datastore.NewKey(ctx, "OAuthAccessToken", teamID, 0, nil)
}

func getBotAccessToken(ctx context.Context, teamID string) (string, error) {
	key := generateOAuthAccessTokenKey(ctx, teamID)
	oauthAccessToken := new(OAuthAccessToken)

	if err := datastore.Get(ctx, key, oauthAccessToken); err == nil {
		return oauthAccessToken.Bot.BotAccessToken, nil
	}

	return "", fmt.Errorf("Bot access token not found for teamID: %s", teamID)
}
