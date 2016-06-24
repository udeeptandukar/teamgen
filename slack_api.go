package teamgen

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type response struct {
	Ok    bool
	Error string
}

func sendMessage(channel string, token string, message string) error {
	v := url.Values{}
	v.Set("token", token)
	v.Set("channel", channel)
	v.Set("text", message)
	v.Set("as_user", "true")

	res, err := http.PostForm("https://slack.com/api/chat.postMessage", v)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody := &response{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(resBody); err != nil {
		return err
	}

	if !resBody.Ok {
		return errors.New(resBody.Error)
	}

	return nil
}
