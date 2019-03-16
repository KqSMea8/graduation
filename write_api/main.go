package write_api

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	log.Info("main init")
	initHttpHandler()
}

func initHttpHandler() {
	log.Info("Init HttpHandler")
	http.HandleFunc("/post", post)
	http.HandleFunc("/update", delete_)
}

