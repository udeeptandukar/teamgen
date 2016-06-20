package teamgen

func addMember(teamID string, channelID string, members []string) slashResponse {
	resp := constructSlashResponse("ephemeral", "Team members added")
	return resp
}
