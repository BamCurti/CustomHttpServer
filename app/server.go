package main

import (
	"log"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/app/server"
)

func main() {
	network := "tcp"
	ip := "0.0.0.0"
	port := "4221"
	s, err := server.New(network, ip, port)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	s.Listen()
}
