package main

import (
	"flag"

	"github.com/lpaarup/img-rftp/pkg/server"
	log "github.com/sirupsen/logrus"
)

var (
	port = flag.Int("port", 69, "port to listen to")
)

func main() {
	flag.Parse()
	var s server.Server
	if err := s.ListenAndServe(*port); err != nil {
		log.Fatalf("error when listening and serving: %v", err)
	}
}
