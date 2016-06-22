package teamgen

import (
	"net/http"
)

func init() {
	http.HandleFunc("/", handleCommand)
	http.HandleFunc("/sendMessage", handleSendMessage)
}
