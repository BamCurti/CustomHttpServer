package main

import (
	"flag"
	"log"

	"github.com/codecrafters-io/http-server-starter-go/app/server"
)

func main() {
	network := "tcp"
	ip := "0.0.0.0"
	port := "4221"

	// enable directory flag
	dirFlag := flag.String("directory", ".", "Directory to serve files from")
	flag.Parse()
	server.DirFlag = dirFlag
	s, err := server.New(network, ip, port)

	if err != nil {
		log.Fatalln(err)
	}

	s.Listen()
}
