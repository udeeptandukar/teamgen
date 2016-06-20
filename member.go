package teamgen

func addMember(teamID string, channelID string, members []string) slashResponse {
	resp := slashResponse{
		ResponseType: "ephemeral",
		Text:         "Team members added",
	}

	return resp
}
