package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/gen1us2k/cloudnative_todo_list/config"
	"github.com/gen1us2k/cloudnative_todo_list/server"
)

func main() {
	c, err := config.Parse()
	if err != nil {
		log.Fatal(err)
	}
	s, err := server.NewServer(c)
	if err != nil {
		log.Fatal(err)
	}
	s.Start()
	log.Fatal(s.Wait())
}
