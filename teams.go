package teamgen

import (
	"math/rand"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// Pair stores pair members
type Pair struct {
	First  string
	Second string
}

// Teams to store information required to generate random team
type Teams struct {
	SlackTeamID        string
	SlackChannelID     string
	Members            []string
	Combinations       []Pair
	MemberExclusions   []Pair
	LastGenerated      []Pair
	LastUpdated        time.Time
	EnableAutoGenerate bool
}

func generateTeamsKey(ctx context.Context, teamID string, channelID string) *datastore.Key {
	keyString := teamID + "-" + channelID
	return datastore.NewKey(ctx, "Teams", keyString, 0, nil)
}

func getAllTeams(ctx context.Context) ([]Teams, error) {
	var teams []Teams
	if _, err := datastore.NewQuery("Teams").GetAll(ctx, &teams); err != nil {
		return teams, err
	}
	return teams, nil
}

func addMember(ctx context.Context, teamID string, channelID string, members []string) SlackCmdResponse {
	key := generateTeamsKey(ctx, teamID, channelID)
	teams := new(Teams)

	if err := datastore.Get(ctx, key, teams); err != nil {
		teams.SlackTeamID = teamID
		teams.SlackChannelID = channelID
		teams.Members = members
		teams.Combinations = generateCombinations(members)
		teams.LastUpdated = time.Now()
	} else {
		teams.Members = members
		teams.Combinations = generateCombinations(members)
		teams.LastUpdated = time.Now()
	}

	if len(teams.Members) > 2 {
		teams.EnableAutoGenerate = true
	} else {
		teams.EnableAutoGenerate = false
	}

	if _, err := datastore.Put(ctx, key, teams); err != nil {
		log.Errorf(ctx, "Error on sending message: %s", err)
		return constructSlackCmdResponse("ephemeral", "Error occurred while adding members. Please try again.")
	}
	return constructSlackCmdResponse("ephemeral", "Team members added: "+strings.Join(members, ", "))
}

func addMemberExclusions(ctx context.Context, teamID string, channelID string, csvPairs []string) SlackCmdResponse {
	key := generateTeamsKey(ctx, teamID, channelID)
	teams := new(Teams)

	exclusionPairs := convertToMemberExclusionPairs(csvPairs)
	if err := datastore.Get(ctx, key, teams); err != nil {
		teams.SlackTeamID = teamID
		teams.SlackChannelID = channelID
		teams.MemberExclusions = exclusionPairs
		teams.LastUpdated = time.Now()
	} else {
		teams.MemberExclusions = exclusionPairs
		teams.LastUpdated = time.Now()
	}

	if _, err := datastore.Put(ctx, key, teams); err != nil {
		log.Errorf(ctx, "Error on sending message: %s", err)
		return constructSlackCmdResponse("ephemeral", "Error occurred while adding members exclusions. Please try again.")
	}
	return constructSlackCmdResponse("ephemeral", "Team members exclusions added: "+strings.Join(csvPairs, " | "))
}

func convertToMemberExclusionPairs(csvPairs []string) []Pair {
	pairs := []Pair{}
	for i := 0; i < len(csvPairs); i++ {
		pair := Pair{First: "", Second: ""}
		splits := strings.Split(csvPairs[i], ",")
		if len(splits) > 0 {
			pair.First = strings.Trim(splits[0], " ")
			if len(splits) > 1 {
				pair.Second = strings.Trim(splits[1], " ")
			}
		}
		pairs = append(pairs, pair)
	}

	return pairs
}

func showConfig(ctx context.Context, teamID string, channelID string) SlackCmdResponse {
	key := generateTeamsKey(ctx, teamID, channelID)
	teams := new(Teams)

	if err := datastore.Get(ctx, key, teams); err == nil {
		text := "Team members: " + strings.Join(teams.Members, ", ")
		text += "\nTeam members exclusions: " + strings.Join(getPairsCSV(teams.MemberExclusions), " | ")
		return constructSlackCmdResponse("ephemeral", text)
	}

	return constructSlackCmdResponse("ephemeral", "No config found.")
}

func toggleAutoGeneration(ctx context.Context, teamID string, channelID string, enable bool) SlackCmdResponse {
	key := generateTeamsKey(ctx, teamID, channelID)
	teams := new(Teams)

	if err := datastore.Get(ctx, key, teams); err != nil {
		teams.SlackTeamID = teamID
		teams.SlackChannelID = channelID
		teams.EnableAutoGenerate = enable
		teams.LastUpdated = time.Now()
	} else {
		teams.EnableAutoGenerate = enable
		teams.LastUpdated = time.Now()
	}

	if _, err := datastore.Put(ctx, key, teams); err != nil {
		log.Errorf(ctx, "Error on enabling auto generation: %s", err)
		return constructSlackCmdResponse("ephemeral", "Action failed. Please try again.")
	}
	return constructSlackCmdResponse("ephemeral", "Action completed")
}

func getRandomTeams(ctx context.Context, teamID string, channelID string) ([]string, error) {
	randomPairs := []string{}
	key := generateTeamsKey(ctx, teamID, channelID)
	teams := new(Teams)

	if err := datastore.Get(ctx, key, teams); err != nil {
		return randomPairs, err
	}
	if len(teams.Combinations) == 0 {
		return randomPairs, nil
	}

	pairs := getRandomPairs(teams.Members, teams.Combinations, teams.MemberExclusions, teams.LastGenerated)
	teams.LastGenerated = pairs

	if _, err := datastore.Put(ctx, key, teams); err != nil {
		return randomPairs, err
	}

	randomPairs = getPairsCSV(pairs)
	return randomPairs, nil
}

func getPairsCSV(pairs []Pair) []string {
	csvPairs := []string{}
	pair := ""
	for i := 0; i < len(pairs); i++ {
		pair = pairs[i].First
		if pairs[i].Second != "" {
			pair = pair + ", " + pairs[i].Second
		}
		csvPairs = append(csvPairs, pair)
	}
	return csvPairs
}

func getRandomPairs(members []string, combinations []Pair, memberExclusions []Pair, lastGenerated []Pair) []Pair {
	pairs := []Pair{}
	pair := Pair{}
	excludes := append(memberExclusions, lastGenerated...)
	combinations = pairSubtraction(combinations, excludes)
	teams := len(members) / 2

	for i := 0; i < teams; i++ {
		if len(combinations) > 0 {
			pair = getWeightedRandomPair(combinations)
			combinations = removePairs(combinations, pair)
			pairs = append(pairs, pair)
			members = removeMembers(members, pair)
		}
	}

	// Get pair from left over members
	if len(members) > 0 {
		pair = Pair{}
		pair.First = members[0]
		if len(members) > 1 {
			pair.Second = members[1]
		} else {
			pair.Second = ""
		}
		pairs = append(pairs, pair)
	}

	return pairs
}

func getWeightedRandomPair(combinations []Pair) Pair {
	minWeightedPairs := []Pair{}
	minWeight := 99999

	for i := 0; i < len(combinations); i++ {
		myWeight := 0
		for j := 0; j < len(combinations); j++ {
			if combinations[i].First == combinations[j].First || combinations[i].First == combinations[j].Second ||
				combinations[i].Second == combinations[j].First || combinations[i].Second == combinations[j].Second {
				myWeight = myWeight + 1
			}
		}
		if myWeight < minWeight {
			minWeight = myWeight
			minWeightedPairs = []Pair{}
			minWeightedPairs = append(minWeightedPairs, combinations[i])
		} else if myWeight == minWeight {
			minWeightedPairs = append(minWeightedPairs, combinations[i])
		}
	}
	return getRandomPair(minWeightedPairs)
}

func getRandomPair(combinations []Pair) Pair {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := r.Intn(len(combinations))
	return combinations[i]
}

func removePairs(combinations []Pair, pair Pair) []Pair {
	pairs := []Pair{}
	for i := 0; i < len(combinations); i++ {
		if pair.First != combinations[i].First && pair.First != combinations[i].Second &&
			pair.Second != combinations[i].First && pair.Second != combinations[i].Second {
			pairs = append(pairs, combinations[i])
		}
	}
	return pairs
}

func removeMembers(members []string, pair Pair) []string {
	newMembers := []string{}
	for i := 0; i < len(members); i++ {
		if pair.First != members[i] && pair.Second != members[i] {
			newMembers = append(newMembers, members[i])
		}
	}
	return newMembers
}

func pairSubtraction(combinations []Pair, excludes []Pair) []Pair {
	pairs := []Pair{}
	for i := 0; i < len(combinations); i++ {
		if !pairExists(combinations[i], excludes) {
			pairs = append(pairs, combinations[i])
		}
	}
	return pairs
}

func getPair(members []string) (Pair, []string) {
	pair := Pair{}
	pair.First, members = popRandomMember(members)
	pair.Second, members = popRandomMember(members)
	return pair, members
}

func popRandomMember(members []string) (string, []string) {
	if len(members) == 0 {
		return "", []string{}
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := r.Intn(len(members))
	v := members[i]
	members[len(members)-1], members[i] = members[i], members[len(members)-1]
	return v, members[:len(members)-1]
}

func pairExists(pair Pair, pairs []Pair) bool {
	for i := 0; i < len(pairs); i++ {
		if (pair.First == pairs[i].First || pair.First == pairs[i].Second) &&
			(pair.Second == pairs[i].First || pair.Second == pairs[i].Second) {
			return true
		}
	}
	return false
}

func generateCombinations(members []string) []Pair {
	pairs := []Pair{}
	n := len(members)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			pair := Pair{First: members[i], Second: members[j]}
			pairs = append(pairs, pair)
		}
	}
	return pairs
}
