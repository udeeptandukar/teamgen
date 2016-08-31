package teamgen

import (
	"net/http"
)

func init() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.Handle("/", http.FileServer(http.Dir("./template/")))
	http.HandleFunc("/cmd", handleCommand)
	http.HandleFunc("/sendMsg", handleSendMessage)
	http.HandleFunc("/oauth", handleOauth)
	http.HandleFunc("/cron/scheduling", handleScheduling)
}
