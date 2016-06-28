package teamgen

import (
	"net/http"
)

func init() {
	http.HandleFunc("/cmd", handleCommand)
	http.HandleFunc("/sendMsg", handleSendMessage)
	http.HandleFunc("/oauth", handleOauth)
	http.HandleFunc("/cron/scheduling", handleScheduling)
}
