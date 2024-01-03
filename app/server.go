package main

import (
	"flag"
	"log"

	"github.com/codecrafters-io/http-server-starter-go/app/server"
)

func main() {
	ip := "0.0.0.0"
	port := "4221"

	// enable directory flag
	dirFlag := flag.String("directory", ".", "Directory to serve files from")
	flag.Parse()
	server.DirFlag = dirFlag

	s, err := server.New(ip, port)

	if err != nil {
		log.Fatalln(err)
	}

	s.Get("/", helloServer)
	s.Get("/echo/{word}", echoHandler)

	s.Serve()
}

func helloServer(req *server.HttpRequest, res *server.HttpResponse) {
	res.ContentType = server.TEXT_PLAIN
	res.StatusCode = server.OK_MSG
	res.Write("")
}

func echoHandler(req *server.HttpRequest, res *server.HttpResponse) {

}
