package teamgen

// TeamGenerator to store information required to generate random team
type TeamGenerator struct {
	Members        []string
	Schedules      []string
	NumberOfTeams  int16
	RandomName     bool
	SlackTeamID    string
	SlackChannelID string
}
type slashResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}
