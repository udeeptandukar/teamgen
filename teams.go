package teamgen

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// Teams to store information required to generate random team
type Teams struct {
	SlackTeamID        string
	SlackChannelID     string
	Members            []string
	Schedules          []string
	NumberOfTeams      int
	RandomName         bool
	MemberCombinations []string
	LastUpdated        time.Time
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
		teams.MemberCombinations = getCombinations(members)
		teams.LastUpdated = time.Now()
	} else {
		teams.Members = members
		teams.NumberOfTeams = len(members) / 2
		teams.MemberCombinations = getCombinations(members)
		teams.LastUpdated = time.Now()
	}

	if _, err := datastore.Put(ctx, key, teams); err != nil {
		log.Errorf(ctx, "Error on sending message: %s", err)
		return constructSlackCmdResponse("ephemeral", "Error occurred while adding members. Please try again.")
	}
	return constructSlackCmdResponse("ephemeral", "Team members added: "+strings.Join(members, ", "))
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

func getRandomMembersMessage(ctx context.Context, teamID string, channelID string) (string, error) {
	result := ""
	key := generateTeamsKey(ctx, teamID, channelID)
	teams := new(Teams)

	if err := datastore.Get(ctx, key, teams); err != nil {
		return "", err
	}
	memberCombinations := teams.MemberCombinations
	index := rand.Intn(len(memberCombinations))
	members := memberCombinations[index]
	memberCombinations = append(memberCombinations[:index], memberCombinations[index+1:]...)
	if len(memberCombinations) == 0 {
		teams.MemberCombinations = getCombinations(teams.Members)
	} else {
		teams.MemberCombinations = memberCombinations
	}
	if _, err := datastore.Put(ctx, key, teams); err != nil {
		return "", err
	}

	result = buildPostMessage(members, teams.NumberOfTeams)
	return result, nil
}

func buildPostMessage(members string, numberOfTeams int) string {
	result := ""
	randomMembers := strings.Split(members, ",")
	teamSize := len(randomMembers) / numberOfTeams
	i := 0
	for i = 0; i < numberOfTeams-1; i++ {
		if i != 0 {
			result = result + "\n"
		}
		result = result + "Team " + strconv.Itoa(i+1) + ": "
		result = result + strings.Join(randomMembers[:(teamSize*(i+1))], ", ")
	}
	result = result + "\nTeam " + strconv.Itoa(i+1) + ": "
	result = result + strings.Join(randomMembers[(teamSize*i):], ", ")
	return result
}

func swapElement(list *[]string, i int, j int) {
	tmp := (*list)[i]
	(*list)[i] = (*list)[j]
	(*list)[j] = tmp
}

func getCombinations(members []string) []string {
	combinations := [][]string{}
	total := len(members)
	curMembers := []string{}
	for i := 0; i < total-1; i++ {
		curMembers = append([]string{}, members...)
		swapElement(&curMembers, 0, i)
		swapElement(&curMembers, 1, i+1)
		combinations = append(combinations, curMembers)

		for j := i + 2; j < total; j++ {
			curMembers = append([]string{}, curMembers...)
			swapElement(&curMembers, 1, j)
			combinations = append(combinations, curMembers)
		}
	}

	results := []string{}
	for i := 0; i < len(combinations); i++ {
		results = append(results, strings.Join(combinations[i], ","))
	}

	return results
}