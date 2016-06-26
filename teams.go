package teamgen

import (
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// Teams to store information required to generate random team
type Teams struct {
	Members        []string
	Schedules      []string
	NumberOfTeams  int
	RandomName     bool
	SlackTeamID    string
	SlackChannelID string
	LastUpdated    time.Time
}

func generateTeamsKey(ctx context.Context, teamID string, channelID string) *datastore.Key {
	keyString := teamID + "-" + channelID
	return datastore.NewKey(ctx, "Teams", keyString, 0, nil)
}

func addMember(ctx context.Context, teamID string, channelID string, members []string) SlackCmdResponse {
	key := generateTeamsKey(ctx, teamID, channelID)
	teams := new(Teams)

	if err := datastore.Get(ctx, key, teams); err != nil {
		teams.SlackTeamID = teamID
		teams.SlackChannelID = channelID
		teams.Members = members
		teams.RandomName = false
		teams.NumberOfTeams = len(members) / 2
		teams.LastUpdated = time.Now()
	} else {
		teams.Members = members
		teams.LastUpdated = time.Now()
	}

	if _, err := datastore.Put(ctx, key, teams); err == nil {
		return constructSlackCmdResponse("ephemeral", "Team members added: "+strings.Join(members, ", "))
	}

	return constructSlackCmdResponse("ephemeral", "Error occurred while adding members. Please try again.")
}

func showConfig(ctx context.Context, teamID string, channelID string) SlackCmdResponse {
	key := generateTeamsKey(ctx, teamID, channelID)
	teams := new(Teams)

	if err := datastore.Get(ctx, key, teams); err == nil {
		text := "Team members: " + strings.Join(teams.Members, ", ")
		text += "\nNo. of teams: " + strconv.Itoa(teams.NumberOfTeams)
		return constructSlackCmdResponse("ephemeral", text)
	}

	return constructSlackCmdResponse("ephemeral", "No config found.")
}

// RandomTeam contains random members for a random team
type RandomTeam struct {
	Name    string
	Members []string
}

func getRandomMembers(ctx context.Context, teamID string, channelID string) []RandomTeam {
	randomMembers := []RandomTeam{}
	key := generateTeamsKey(ctx, teamID, channelID)
	teams := new(Teams)

	if err := datastore.Get(ctx, key, teams); err != nil {
		//TODO: Create random teams
	}
	return randomMembers
}
